#pragma once

#include <string>
#include <vector>

#include "src/carnot/exec/exec_state.h"
#include "src/carnot/exec/row_batch.h"
#include "src/carnot/exec/row_descriptor.h"
#include "src/carnot/plan/operators.h"
#include "src/common/common.h"

namespace pl {
namespace carnot {
namespace exec {

enum class ExecNodeType : int8_t {
  kSourceNode = 0,
  kSinkNode = 1,
  kProcessingNode = 2,
};

/**
 * This is the base class for the execution nodes in Carnot.
 */
class ExecNode {
 public:
  ExecNode() = delete;
  virtual ~ExecNode() = default;

  /**
   * Init is called with plan & schema information.
   * @param plan_node the plan class of the node.
   * @param output_descriptor The output column schema of row batches.
   * @param input_descriptors The input column schema of row batches.
   * @return
   */
  Status Init(const plan::Operator& plan_node, const RowDescriptor& output_descriptor,
              std::vector<RowDescriptor> input_descriptors) {
    return InitImpl(plan_node, output_descriptor, input_descriptors);
  }

  /**
   * Setup internal data structures, perform validation, etc.
   * @param exec_state The execution state.
   * @return The status of the prepare.
   */
  Status Prepare(ExecState* exec_state) { return PrepareImpl(exec_state); }

  /**
   * Acquire memory resources, etc.
   * @param exec_state The execution state.
   * @return
   */
  Status Open(ExecState* exec_state) { return OpenImpl(exec_state); }

  /**
   * Close is where cleanup should take place. This includes cleaning up objects.
   * It is highly recomended that a default destructor be used and cleanup peformed here,
   * since at the end a query the data is batch deleted and ordering is not guaranteed.
   * @param exec_state The execution state.
   * @return The status of the Finalize.axs
   */
  Status Close(ExecState* exec_state) { return CloseImpl(exec_state); }

  /**
   * GenerateNext is called to produce the next row batch. This is only valid
   * on source nodes (and will result in an error on other nodes).
   * @param exec_state The execution state.
   * @return The status of the execution.
   */
  Status GenerateNext(ExecState* exec_state) {
    DCHECK(type() == ExecNodeType::kSourceNode);
    return GenerateNextImpl(exec_state);
  }

  /**
   * Consume the next row batch. This function is only valid for Sink and Processing
   * Nodes.
   *
   * This needs to be careful to forward the output batch to all children.
   *
   * @param exec_state The execution state.
   * @param rb The input row batch.
   * @return The Status of consumption.
   */
  Status ConsumeNext(ExecState* exec_state, const RowBatch& rb) {
    DCHECK(type() == ExecNodeType::kSinkNode || type() == ExecNodeType::kProcessingNode);

    return ConsumeNextImpl(exec_state, rb);
  }

  /**
   * Check if it's a source node.
   */
  bool IsSource() { return type() == ExecNodeType::kSourceNode; }

  /**
   * Check if it's a sink node.
   */
  bool IsSink() { return type() == ExecNodeType::kSinkNode; }

  /**
   * Check if it's a processing node.
   */
  bool IsProcessing() { return type() == ExecNodeType::kProcessingNode; }

  /**
   * Get a debug string for the node.
   * @return the debug string/
   */
  std::string DebugString() { return DebugStringImpl(); }

  /**
   * Add a new child node where data is forwarded.
   * This node will not own the child. The lifetime of the child should
   * exceed the lifetime of this node.
   *
   * @param child Another execution node.
   */
  void AddChild(ExecNode* child) { children_.emplace_back(child); }

  /**
   * Get the type of the execution node.
   * @return the ExecNodeType.
   */
  ExecNodeType type() { return type_; }

  /**
   * @ return the children of the execution node.
   */
  std::vector<ExecNode*> children() { return children_; }

 protected:
  /**
   * Send data to children row batches.
   * @param exec_state The exec state.
   * @param rb The row batch to send.
   * @return Status of children execution.
   */
  Status SendRowBatchToChildren(ExecState* exec_state, const RowBatch& rb) {
    for (auto* child : children_) {
      PL_RETURN_IF_ERROR(child->ConsumeNext(exec_state, rb));
    }
    return Status::OK();
  }

 protected:
  explicit ExecNode(ExecNodeType type) : type_(type) {}

  // Defines the protected implementations of the non-virtual interface functions
  // defined above.
  virtual std::string DebugStringImpl() = 0;
  virtual Status InitImpl(const plan::Operator& plan_node, const RowDescriptor& output_descriptor,
                          const std::vector<RowDescriptor>& input_descriptors) = 0;
  virtual Status PrepareImpl(ExecState* exec_state) = 0;
  virtual Status OpenImpl(ExecState* exec_state) = 0;
  virtual Status CloseImpl(ExecState* exec_state) = 0;

  virtual Status GenerateNextImpl(ExecState*) {
    return error::Unimplemented("Implement in derived class (if source)");
  }

  virtual Status ConsumeNextImpl(ExecState*, const RowBatch&) {
    return error::Unimplemented("Implement in derived class (if sink or processing)");
  }

  bool is_closed() { return is_closed_; }

 private:
  // Unowned reference to the children. Must remain valid for the duration of query.
  std::vector<ExecNode*> children_;
  bool is_closed_ = false;
  // The type of execution node.
  ExecNodeType type_;
};

/**
 * Processing node is the base class for anything that computes
 * producing 1:1 or N:M records. For example: Agg, Map, etc.
 */
class ProcessingNode : public ExecNode {
 public:
  ProcessingNode() : ExecNode(ExecNodeType::kProcessingNode) {}
};

/**
 * Source node is the base class for anything that produces records from some source.
 * For example: MemorySource.
 */
class SourceNode : public ExecNode {
 public:
  SourceNode() : ExecNode(ExecNodeType::kSourceNode) {}
  virtual bool BatchesRemaining() = 0;
  virtual bool NextBatchReady() = 0;
};

/**
 * Sink node is the base class for anything that consumes records and writes to some sink.
 * For example: MemorySink.
 */
class SinkNode : public ExecNode {
 public:
  SinkNode() : ExecNode(ExecNodeType::kSinkNode) {}
};

}  // namespace exec
}  // namespace carnot
}  // namespace pl

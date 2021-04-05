#include <gmock/gmock.h>
#include <google/protobuf/text_format.h>
#include <gtest/gtest.h>

#include <unordered_map>
#include <utility>
#include <vector>

#include <pypa/parser/parser.hh>

#include "src/carnot/planner/compiler/optimizer/limit_push_down.h"
#include "src/carnot/planner/compiler/test_utils.h"

namespace pl {
namespace carnot {
namespace planner {
namespace compiler {

using table_store::schema::Relation;
using ::testing::_;

using ::testing::Contains;
using ::testing::ContainsRegex;
using ::testing::ElementsAre;
using ::testing::ElementsAreArray;
using ::testing::Not;
using ::testing::UnorderedElementsAre;

TEST_F(OperatorTests, simple_no_op) {
  Relation relation({types::DataType::INT64, types::DataType::INT64}, {"abc", "xyz"});
  MemorySourceIR* src = MakeMemSource(relation);
  FilterIR* filter = MakeFilter(src, MakeEqualsFunc(MakeColumn("abc", 0), MakeColumn("xyz", 0)));
  LimitIR* limit = MakeLimit(filter, 10);
  MakeMemSink(limit, "foo", {});

  LimitPushdownRule rule;
  auto result = rule.Execute(graph.get());
  ASSERT_OK(result);
  EXPECT_FALSE(result.ValueOrDie());
}

TEST_F(OperatorTests, simple) {
  Relation relation({types::DataType::INT64, types::DataType::INT64}, {"abc", "xyz"});

  MemorySourceIR* src = MakeMemSource(relation);
  MapIR* map1 = MakeMap(src, {{"abc", MakeColumn("abc", 0)}}, false);
  MapIR* map2 = MakeMap(map1, {{"def", MakeColumn("abc", 0)}}, false);
  LimitIR* limit = MakeLimit(map2, 10);
  MemorySinkIR* sink = MakeMemSink(limit, "foo", {});

  LimitPushdownRule rule;
  auto result = rule.Execute(graph.get());
  ASSERT_OK(result);
  EXPECT_TRUE(result.ValueOrDie());

  EXPECT_THAT(sink->parents(), ElementsAre(map2));
  EXPECT_THAT(map2->parents(), ElementsAre(map1));
  EXPECT_EQ(1, map1->parents().size());
  auto new_limit = map1->parents()[0];
  EXPECT_MATCH(new_limit, Limit());
  EXPECT_EQ(10, static_cast<LimitIR*>(new_limit)->limit_value());
  EXPECT_THAT(new_limit->parents(), ElementsAre(src));

  // Ensure the limit inherits the relation from its parent, not the previous location.
  EXPECT_EQ(new_limit->relation(), relation);

  // Rerunning the rule shouldn't cause the limits to get cloned again.
  result = rule.Execute(graph.get());
  ASSERT_OK(result);
  EXPECT_FALSE(result.ValueOrDie());
}

TEST_F(OperatorTests, multi_branch_union) {
  Relation relation1({types::DataType::INT64, types::DataType::INT64}, {"abc", "xyz"});
  Relation relation2({types::DataType::INT64}, {"abc"});

  MemorySourceIR* src1 = MakeMemSource(relation1);
  MapIR* map1 = MakeMap(src1, {{"abc", MakeColumn("abc", 0)}}, false);

  MemorySourceIR* src2 = MakeMemSource(relation2);
  MapIR* map2 = MakeMap(src2, {{"abc", MakeColumn("abc", 0)}}, false);

  UnionIR* union_node = MakeUnion({map1, map2});
  ASSERT_OK(union_node->SetRelation(relation2));
  MapIR* map3 = MakeMap(union_node, {{"abc", MakeColumn("abc", 0)}}, false);

  LimitIR* limit = MakeLimit(map3, 10);
  MemorySinkIR* sink = MakeMemSink(limit, "foo", {});

  LimitPushdownRule rule;
  auto result = rule.Execute(graph.get());
  ASSERT_OK(result);
  EXPECT_TRUE(result.ValueOrDie());

  EXPECT_THAT(sink->parents(), ElementsAre(map3));

  // There should be 3 copies of the limit -- one before of each branch of the
  // union, and one after the union.
  EXPECT_EQ(1, map3->parents().size());
  auto new_limit = map3->parents()[0];
  EXPECT_EQ(10, static_cast<LimitIR*>(new_limit)->limit_value());
  EXPECT_THAT(new_limit->parents(), ElementsAre(union_node));
  EXPECT_EQ(new_limit->relation(), relation2);

  EXPECT_THAT(union_node->parents(), ElementsAre(map1, map2));

  EXPECT_EQ(1, map1->parents().size());
  auto new_limit1 = map1->parents()[0];
  EXPECT_MATCH(new_limit1, Limit());
  EXPECT_EQ(10, static_cast<LimitIR*>(new_limit1)->limit_value());
  EXPECT_THAT(new_limit1->parents(), ElementsAre(src1));
  // Ensure the limit inherits the relation from its parent, not the previous location.
  EXPECT_EQ(new_limit1->relation(), relation1);

  EXPECT_EQ(1, map2->parents().size());
  auto new_limit2 = map2->parents()[0];
  EXPECT_MATCH(new_limit2, Limit());
  EXPECT_EQ(10, static_cast<LimitIR*>(new_limit2)->limit_value());
  EXPECT_THAT(new_limit2->parents(), ElementsAre(src2));
  // Ensure the limit inherits the relation from its parent, not the previous location.
  EXPECT_EQ(new_limit2->relation(), relation2);

  // Rerunning the rule shouldn't cause the limits to get cloned again.
  result = rule.Execute(graph.get());
  ASSERT_OK(result);
  EXPECT_FALSE(result.ValueOrDie());
}

}  // namespace compiler
}  // namespace planner
}  // namespace carnot
}  // namespace pl
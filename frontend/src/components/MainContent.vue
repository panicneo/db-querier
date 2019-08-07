<template>
  <el-table
    :data="queryResult"
    border
    stripe
    height="900"
    @row-dblclick="handleClick"
  >
    <el-table-column v-if="queryResult.length != 0" type="expand">
      <template slot-scope="props">
        <el-container style="overflow-x: auto;">
          <pre v-highlightjs><code class="json">{{props.row}} </code></pre>
        </el-container>
      </template>
    </el-table-column>
    <el-table-column
      v-if="queryResult.length != 0"
      type="index"
    ></el-table-column>
    <el-table-column
      v-for="(col, index) in columns"
      :key="index"
      :prop="col"
      :label="col"
      max-width="300"
    ></el-table-column>
  </el-table>
</template>

<script>
import _ from "lodash/core";
import { Table, TableColumn, Container } from "element-ui";
export default {
  components: {
    ElTable: Table,
    ElTableColumn: TableColumn,
    ElContainer: Container
  },
  computed: {
    queryResult: function() {
      return this.$store.state.queryResult;
    },
    columns: function() {
      return _.keys(this.queryResult[0]);
    }
  },
  methods: {
    handleClick(row) {
      let vm = this;
      this.$copyText(JSON.stringify(row)).then(
        () => vm.$notify.success("复制成功"),
        () => vm.$notify.error("复制失败")
      );
    }
  }
};
</script>

<style scope>
.el-table .cell {
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>

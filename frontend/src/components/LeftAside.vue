<template>
  <el-menu :default-openeds="['1']">
    <el-menu-item
      v-for="(name, index) in selections"
      :key="index"
      :index="index.toString()"
      @click="onClickQueryName(name)"
      >{{ name }}</el-menu-item
    >
  </el-menu>
</template>

<script>
import _ from "lodash/core";
import { Menu, MenuItem } from "element-ui";
export default {
  components: {
    ElMenu: Menu,
    ElMenuItem: MenuItem
  },
  computed: {
    selections: function() {
      return _.map(this.$store.state.queries, "name");
    }
  },
  mounted() {
    this.fetchQueries();
  },
  methods: {
    onClickQueryName: function(name) {
      this.$store.commit("updateCurrent", { current: name });
    },
    fetchQueries: async function() {
      let resp = await this.$axios.get("/api/query");
      this.$store.commit("updateQueries", { queries: resp.data });
    }
  }
};
</script>

<style></style>

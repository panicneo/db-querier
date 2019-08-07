<template>
  <el-form
    v-if="isSelected"
    ref="kvMap"
    :model="kvMap"
    label-position="right"
    label-width="auto"
    status-icon
    :rules="paramValidators"
  >
    <el-form-item>
      <el-button type="primary" :loading="loading" @click="onSubmit('kvMap')"
        >查询</el-button
      >
      <el-button type="primary" :loading="loading" @click="onCopyAll"
        >复制All</el-button
      >
    </el-form-item>
    <el-form-item
      v-for="(param, index) in params"
      :key="index"
      :label="param.name"
      :prop="param.name"
    >
      <el-autocomplete
        v-if="param.type === 'number'"
        v-model.number="kvMap[param.name]"
        :fetch-suggestions="(qs, cb) => fetchHistory(qs, cb, param.name)"
      ></el-autocomplete>
      <el-input v-else v-model="kvMap[param.name]"></el-input>
    </el-form-item>
  </el-form>
</template>

<script>
import { Form, FormItem, Autocomplete, Button } from "element-ui";
import _ from "lodash/core";

export default {
  components: {
    ElForm: Form,
    ElFormItem: FormItem,
    ElAutocomplete: Autocomplete,
    ElButton: Button
  },
  data() {
    return {
      kvMap: {},
      loading: false
    };
  },
  computed: {
    params: function() {
      return this.$store.state.current.params;
    },
    isSelected: function() {
      return this.$store.state.current != "";
    },
    paramValidators: function() {
      return this.$store.state.validators[this.$store.state.current.name];
    }
  },
  methods: {
    curHistories: async function(field) {
      let histories = await window.localStorage.getItem(
        `query-his:${this.$store.state.current.name}:${field}`
      );
      if (!histories) {
        histories = [];
      } else {
        histories = JSON.parse(histories);
      }
      return histories;
    },
    onSubmit: function(formName) {
      this.loading = true;
      this.$refs[formName].validate(valid => {
        if (!valid) {
          this.loading = false;
          return false;
        } else {
          const name = this.$store.state.current.name;
          let url = `/api/query/${name}`;
          this.$axios
            .get(url, { params: this.kvMap })
            .then(resp => {
              if (resp.data.length === 0) {
                this.$notify.warning("查询结果为空");
                return;
              }
              this.$store.commit("updateQueryResult", {
                queryResult: resp.data
              });
            })
            .finally(() => {
              this.loading = false;
              this.saveHistory();
            });
        }
      });
    },
    onCopyAll: function() {
      let vm = this;
      let data = JSON.stringify(this.$store.state.queryResult);
      this.$copyText(data).then(
        () => vm.$notify.success("复制成功"),
        () => vm.$notify.error("复制失败")
      );
    },
    saveHistory: async function() {
      _.forEach(this.kvMap, async (v, k) => {
        let his = await this.curHistories(k);
        his.unshift(v);
        his = Array.from(new Set(his));
        await window.localStorage.setItem(
          `query-his:${this.$store.state.current.name}:${k}`,
          JSON.stringify(his.slice(0, 9))
        );
      });
    },
    fetchHistory: async function(queryString, cb, name) {
      let histories = await this.curHistories(name);
      queryString &&
        histories.filter(his => String(his).startsWith(queryString));
      let result = [];
      _.forEach(histories, x => result.push({ value: String(x) }));
      cb(result);
    }
  }
};
</script>

<style scoped>
.el-form-item {
  width: 500px;
}
</style>

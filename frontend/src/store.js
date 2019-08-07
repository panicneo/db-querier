import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

import _ from "lodash/core";
export default new Vuex.Store({
  state: { queries: {}, current: "", validators: {}, queryResult: [] },
  mutations: {
    updateQueries: (state, payload) => {
      state.queries = payload.queries;
      _.forEach(payload.queries, query => {
        state.queries[query.name] = query;
        let validator = {};
        _.forEach(query.params, param => {
          validator[param.name] = [
            {
              type: param.type,
              message: `${param.name}字段是一个 ${param.type}类型`,
              trigger: "blur"
            },
            {
              required: true,
              message: `${param.name}不能为空`,
              trigger: "blur"
            }
          ];
        });
        state.validators[query.name] = validator;
      });
    },
    updateCurrent: (state, payload) =>
      (state.current = state.queries[payload.current]),
    updateQueryResult: (state, payload) => {
      state.queryResult = payload.queryResult;
    }
  },
  actions: {}
});

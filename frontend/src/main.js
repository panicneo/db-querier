import "./plugins/element.js";
import "./plugins/axios";
import "./plugins/clipboard";
import "./plugins/highlight";

import Vue from "vue";

import App from "./App.vue";
import store from "./store";

// Vue.config.productionTip = false;

new Vue({ store, render: h => h(App) }).$mount("#app");

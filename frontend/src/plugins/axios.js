"use strict";

import Vue from "vue";
import axios from "axios";

const isProd = process.env.NODE_ENV === "production";
let config = {
  baseURL: isProd ? "" : "http://localhost:5000",
  timeout: 60 * 1000, // Timeout
  withCredentials: isProd // Check cross-site Access-Control
};

const _axios = axios.create(config);
import { Notification } from "element-ui";

_axios.interceptors.request.use(
  function(config) {
    // Do something before request is sent
    return config;
  },
  function(error) {
    // Do something with request error
    Notification.error(error.response.data.error || "服务器连接失败");
    return Promise.reject(error.response);
  }
);

// Add a response interceptor
_axios.interceptors.response.use(
  function(response) {
    // Do something with response data
    return response;
  },
  function(error) {
    // Do something with response error
    Notification.error(
      error.response ? error.response.data.error : "服务器发生错误"
    );
    return Promise.reject(error.response);
  }
);

Plugin.install = function(Vue) {
  Vue.axios = _axios;
  window.axios = _axios;
  Object.defineProperties(Vue.prototype, {
    axios: {
      get() {
        return _axios;
      }
    },
    $axios: {
      get() {
        return _axios;
      }
    }
  });
};

Vue.use(Plugin);

export default Plugin;

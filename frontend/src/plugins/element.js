import { Notification } from "element-ui";
import Vue from "vue";

Vue.component(Notification);
Vue.prototype.$notify = Notification;

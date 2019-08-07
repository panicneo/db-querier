import { AxiosInstance } from "axios";
import { ElNotification } from "element-ui/types/notification";
import Vue from "vue";

declare module "vue/types/vue" {
  interface Vue {
    $axios: AxiosInstance;
    $notify: ElNotification;
  }
}

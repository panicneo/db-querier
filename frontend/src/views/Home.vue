<template>
  <el-container id="container">
    <left-aside />
    <el-main direction="vertical">
      <main-header />
      <main-content />
    </el-main>
  </el-container>
</template>

<script>
const LeftAside = () => import("../components/LeftAside.vue");
const MainContent = () => import("../components/MainContent.vue");
const MainHeader = () => import("../components/MainHeader.vue");
import { Container, Main } from "element-ui";
export default {
  components: {
    ElContainer: Container,
    ElMain: Main,
    LeftAside,
    MainContent,
    MainHeader
  }
};
</script>

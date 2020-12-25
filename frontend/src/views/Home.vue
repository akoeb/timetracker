<template>
  <div class="home">
    <img alt="Vue logo" src="../assets/logo.png" />
    <!-- HelloWorld msg="Welcome to Your Vue.js App" / -->
    <ProjectList></ProjectList>
  </div>
</template>

<script>
// @ is an alias to /src
// import HelloWorld from "@/components/HelloWorld.vue";
import ProjectList from "@/components/ProjectList.vue";

export default {
  name: "Home",
  components: {
    ProjectList
  }
};
</script>

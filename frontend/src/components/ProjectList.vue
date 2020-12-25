<template>
  <div class="projectList">
    <h1>Projects</h1>
    <ul v-if="projects.length">
      <ProjectItem
        v-for="project in projects"
        :key="project.id"
        :project="project"
        @remove="removeProject"
      />
    </ul>
    <p class="none" v-else>
      Nothing left in the list. Add a new project in the input below.
    </p>
    <form @submit.prevent="addProject">
      <input
        type="text"
        name="project-name"
        v-model="newProjectName"
        placeholder="New Project Name"
      />
      <input
        type="text"
        name="project-client-name"
        v-model="newProjectClientName"
        placeholder="New Project Client Name"
      />
      <select v-model="newProjectStatus">
        <option disabled value="">Initial Project State</option>
        <option>OPEN</option>
        <option>CLOSED</option>
        <option>ACTIVE</option>
      </select>
      <input type="submit" />
    </form>
  </div>
</template>

<script>
import { createProject, deleteProject } from "@/utils";
import ProjectItem from "./ProjectItem.vue";
export default {
  name: "ProjectList",
  components: {
    ProjectItem
  },
  data() {
    return {
      projects: []
    };
  },
  methods: {
    addProject() {
      const trimmedName = this.newProjectName.trim();
      const trimmedClientName = this.newProjectClientName.trim();

      if (trimmedName) {
        this.projects.push(
          createProject(trimmedName, trimmedClientName, this.newProjectStatus)
        );
      }
      this.newProjectName = "";
      this.newProjectClientName = "";
      this.newProjectStatus = "";
    },

    removeProject(item) {
      this.projects = this.projects.filter(project => project !== item);
      deleteProject(item.id);
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped></style>

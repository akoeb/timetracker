<template>
  <li>
    <span>{{ project.name }}</span>
    <button @click.prevent="$emit('remove', project)">Remove</button>
  </li>
</template>

<script setup>
import { RouterView } from 'vue-router'
</script>

<template>
  <div class="container-fluid">
    <div class="row">
      <!-- Rimosso il codice della barra laterale -->
      <main class="col-12 px-md-4">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<style>
/* Rimosso lo stile per la barra laterale */
</style>

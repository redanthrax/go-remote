<script lang="ts">
export default {
    data() {
        return {
            agents: [],
        }
    },
    mounted() {
        this.fetchAgents()
    },
    methods: {
        fetchAgents() {
            fetch('http://localhost:8080/agent')
                .then(response => response.json())
                .then(json => {
                    this.agents = json
                })
        }
    }
}
</script>

<template>
  <div class="agents">
    <h1>Agents List</h1>
    <div class="agents">
        <div v-for="agent in agents" :key="agent.ID">
            <div class="agent">
                <a :href="'/agent?id=' + agent.ID" target="_blank">Agent ID: {{agent.ID}} - Ready: {{agent.Ready}}</a>
            </div>
        </div>
    </div>
  </div>
</template>

<style>
@media (min-width: 1024px) {
  .agents {
    min-height: 100vh;
    display: flex;
    align-items: center;
  }
}
</style>

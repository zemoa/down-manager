<script lang="ts">
    import { onMount } from "svelte";
    import { linkStore } from "../business/store/link-store";

    onMount(async () => {
        linkStore.fetchingLinks.set(true)
        fetch("http://localhost:8080/links")
        .then(response => response.json())
        .then(data => {
            console.log(data)
        })
        .finally(() => {
            linkStore.fetchingLinks.set(false)
        })
    })
</script>
<main>
    <h1>Welcome to SvelteKit</h1>
    <p>Visit <a href="https://kit.svelte.dev">kit.svelte.dev</a> to read the documentation Hello</p>
</main>


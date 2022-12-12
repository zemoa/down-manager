<script lang="ts">
    import { onMount } from "svelte";
    import { linkStore } from "../business/store/link-store";
    const {fetchingLinks, links} = linkStore
    onMount(async () => {
        linkStore.retrieveLinks();
    })
</script>
<main>
    <h1>Welcome to SvelteKit</h1>
    <p>Visit <a href="https://kit.svelte.dev">kit.svelte.dev</a> to read the documentation Hello</p>
    <p>Fetching : {$fetchingLinks}</p>
    <ul>
    {#each $links as link}
        <li>{link.Ref}</li>
    {/each}
    </ul>
</main>


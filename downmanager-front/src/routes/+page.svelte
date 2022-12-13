<script lang="ts">
    import { onMount } from "svelte";
    import { linkStore } from "../business/store/link-store";
    import Button, {Label} from "@smui/button";
	import AddDialog from "../components/AddDialog.svelte";
    
    const {fetchingLinks, links} = linkStore
    let openAddDialog = false;
    onMount(async () => {
        linkStore.retrieveLinks();
    })
</script>
<main>
    <p>Fetching : {$fetchingLinks}</p>
    <Button on:click={() => (openAddDialog = true)}>
        <Label>Add Link</Label>
    </Button>
    <AddDialog bind:open={openAddDialog}></AddDialog>
    
    <ul>
    {#each $links as link}
        <li>{link.Ref}</li>
    {/each}
    </ul>
</main>


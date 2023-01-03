<script lang="ts">
    import { onMount } from "svelte";
    import { linkStore } from "../business/store/link-store";
    import { webSocketStore } from "../business/store/websocket-store";
    import Button, {Label} from "@smui/button";
	import AddDialog from "../components/AddDialog.svelte";
    import DeleteDialog from "../components/DeleteDialog.svelte";
	import LinkItem from "../components/LinkItem.svelte";
    
    const {fetchingLinks, links} = linkStore
    let openAddDialog = false;
    let removeDialogData = {
        open: false,
        linkRef: ""
    };
    onMount(async () => {
        webSocketStore.connect();
        linkStore.retrieveLinks();
    })
    
</script>
<main>
    <p>Fetching : {$fetchingLinks}</p>
    <Button on:click={() => (openAddDialog = true)}>
        <Label>Add Link</Label>
    </Button>
    <AddDialog bind:open={openAddDialog}></AddDialog>
    <DeleteDialog bind:open={removeDialogData.open} bind:linkRef={removeDialogData.linkRef}></DeleteDialog>
    
    <ul>
    {#each $links as link}
        <li><LinkItem bind:link bind:removeDialogData></LinkItem></li>
    {/each}
    </ul>
</main>


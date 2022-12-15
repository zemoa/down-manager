<script lang="ts">
    import { onMount } from "svelte";
    import { linkStore } from "../business/store/link-store";
    import Button, {Label} from "@smui/button";
	import AddDialog from "../components/AddDialog.svelte";
    import DeleteDialog from "../components/DeleteDialog.svelte";
    
    const {fetchingLinks, links} = linkStore
    let openAddDialog = false;
    let removeDialogData = {
        open: false,
        linkRef: ""
    };
    onMount(async () => {
        linkStore.retrieveLinks();
    })
    function onRemoveCalled(linkref: string) {
        removeDialogData = {
            open: true,
            linkRef: linkref
        }
    }
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
        <li>{link.Ref} - <Button on:click={() => (onRemoveCalled(link.Ref))}>Remove</Button></li>
    {/each}
    </ul>
</main>


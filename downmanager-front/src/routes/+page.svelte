<script lang="ts">
    import { onMount } from "svelte";
    import { linkStore } from "../business/store/link-store";
    import Button, {Label} from "@smui/button";
    import Dialog, {Title, Content, Actions} from "@smui/dialog"; 
    import Textfield from "@smui/textfield";
    const {fetchingLinks, links} = linkStore
    let openAddDialog = false;
    let link = "";
    onMount(async () => {
        linkStore.retrieveLinks();
    })
</script>
<main>
    <p>Fetching : {$fetchingLinks}</p>
    <Button on:click={() => (openAddDialog = true)}>
        <Label>Add Link</Label>
    </Button>
    <Dialog
        bind:open={openAddDialog}
        aria-labelledby="simple-title"
        aria-describedby="simple-content"
    >

    <Title id="simple-title">Add a Link</Title>
    <Content id="simple-content">
        <div>
            <Textfield bind:value={link} label="Link" variant="outlined">
            </Textfield>
        </div>
    </Content>
    <Actions>
        <Button on:click={() => (console.log('no'))}>
        <Label>No</Label>
        </Button>
        <Button on:click={() => (console.log('yes'))}>
        <Label>Yes</Label>
        </Button>
    </Actions>
    </Dialog>
    <ul>
    {#each $links as link}
        <li>{link.Ref}</li>
    {/each}
    </ul>
</main>


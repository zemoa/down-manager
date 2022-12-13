<script lang="ts">
    import Dialog, {Title, Content, Actions} from "@smui/dialog"; 
    import Textfield from "@smui/textfield";
    import Button, {Label} from "@smui/button";
    import { linkStore } from "../business/store/link-store";

    let link = "";
    export let open: boolean;
    function onAdd() {
        linkStore.addLink(link)
        link = "";
    }
</script>

<Dialog
        bind:open
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
        <Button on:click={() => (link = "")}>
            <Label>No</Label>
        </Button>
        <Button on:click={onAdd}>
            <Label>Yes</Label>
        </Button>
    </Actions>
</Dialog>
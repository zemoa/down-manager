<script lang="ts">
    import Dialog, {Title, Content, Actions} from "@smui/dialog"; 
    import Textfield from "@smui/textfield";
    import Button, {Label} from "@smui/button";
    import { linkStore } from "../business/store/link-store";

    export let linkRef: string;
    export let open: boolean;
    function onRemove() {
        linkStore.removeLink(linkRef)
    }
</script>

<Dialog
        bind:open
        aria-labelledby="simple-title"
        aria-describedby="simple-content"
    >

    <Title id="simple-title">Delete a link</Title>
    <Content id="simple-content">
        Would you delete this link {linkRef}
    </Content>
    <Actions>
        <Button>
            <Label>No</Label>
        </Button>
        <Button on:click={onRemove}>
            <Label>Yes</Label>
        </Button>
    </Actions>
</Dialog>
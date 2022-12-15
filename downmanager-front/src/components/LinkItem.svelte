<script lang="ts">
	import Button from "@smui/button";
	import type { Link } from "src/business/model/link-model";

    export let link:Link
    export let removeDialogData : {
        open: boolean,
        linkRef: string
    }
    function onRemoveCalled(linkref: string) {
        removeDialogData = {
            open: true,
            linkRef: linkref
        }
    }
</script>
{link.Ref} - <Button on:click={() => (onRemoveCalled(link.Ref))}>Remove</Button>
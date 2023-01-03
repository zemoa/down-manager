<script lang="ts">
	import Button, {Icon} from "@smui/button";
	import type { Link, LinkItem } from "../business/model/link-model";
	import { linkStore } from "../business/store/link-store";

    export let link:LinkItem
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
{link.Ref} ({link.Filename}) - state : {link.Running ? "running": "stopped"} - {link.Percent} %
<Button on:click={() => (linkStore.startDownload(link.Ref))}>
    <Icon class="material-icons">play_arrow</Icon>
</Button>
<Button on:click={() => (linkStore.stopDownload(link.Ref))}>
    <Icon class="material-icons">pause</Icon>
</Button> 
<Button on:click={() => (onRemoveCalled(link.Ref))}>Remove</Button>
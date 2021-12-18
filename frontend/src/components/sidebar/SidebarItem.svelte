<script>
import {Icon} from "svelte-awesome";

export let name = undefined;
export let icon = undefined;
export let to = undefined

let hover = false;
</script>

<div class="header-menu-item">
    <div><Icon class="header-menu-item-icon" data={icon} label={name} /></div>
    <div><a href="#{to}">{name}</a></div>
</div>

<style>
    :global(.header-menu-item-icon) {
      vertical-align: middle;
    }

    :global(.sidebar-link:hover) {
        color: purple;
        margin: auto;
    }

    :global(.header-menu-item) {
        padding: .5em;
        display: grid;
        grid-template-columns: [icon] 24px [label] auto;
    }
</style>
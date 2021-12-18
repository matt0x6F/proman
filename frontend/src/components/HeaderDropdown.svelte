<script>
    import {Dropdown, DropdownShell} from "attractions";
    import {Button} from "svelma";
    import {Icon} from "svelte-awesome";
    import {ellipsisV, gears, listAlt} from "svelte-awesome/icons"
    import SidebarItem from "./sidebar/SidebarItem.svelte";
</script>

<div>
    <DropdownShell class="header-dropdown" let:toggle>
        <Button on:click={toggle} type="is-dark" size="is-small">
            <Icon  class="header-menu-icon" data={ellipsisV} label={"Menu"} />
        </Button>
        <Dropdown class="header-menu">
            <div class="header-menu-content">
                <SidebarItem name="Projects" icon={listAlt} to="/" />
                <SidebarItem name="Settings" icon={gears} to="/app/settings" />
            </div>
        </Dropdown>
    </DropdownShell>
</div>

<style lang="scss">
    @use 'theme.scss';

    :global(.header-menu-icon) {
        color: whitesmoke;
        float: right;
        margin: auto;
    }

    /* menu content*/
    :global(.header-menu) {
        position: absolute !important;
        background-color: white;
        display: block;
        float-offset: 1em;
        left: auto !important;
        right: 0 !important;
    }

    /* menu container */
    :global(.header-dropdown) {
        position: relative !important;
        display: inline-block;
    }

    :global(.header-menu-content) {
        padding: 1em;
    }
</style>
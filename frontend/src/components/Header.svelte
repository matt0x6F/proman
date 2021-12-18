<script>
    import HeaderDropdown from "./HeaderDropdown.svelte";
</script>

<div>
    <h1>proman</h1>
    <HeaderDropdown />
</div>

<style>
    div {
        display: grid;
        grid-template-columns: [title] auto [menu] 2em;
        background-color: #3d3d3d;
        width: 100%;
        padding: 1em;
    }

    div h1 {
        color: whitesmoke;
        margin: auto auto auto 0px;
    }
</style>
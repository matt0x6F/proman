<script>
    import DirectorySelector from "../components/settings/DirectorySelector.svelte";
    import {Headline} from "attractions";
    import Editors from "../components/settings/Editors.svelte";

    let warnings = {};
    let config ={};
    let loading = true;
    let error = undefined;

    window.backend.Config.Get().then((data) => {
        config = data;

        window.wails.Events.Emit("validate_config", config);

        loading = false;
    }).catch((err) => {
        error = err;
        loading = false;
    });

    window.wails.Events.On("validate.config.completed", (errors) => warnings = errors);

    // called on blur
    function handleInput(e) {
        config[e.target.name] = e.target.value;

        window.wails.Events.Emit("config.update", config);
        window.wails.Events.Emit("validate.config", config);
    }
</script>

<div>
    <Headline>Settings</Headline>
    {#if loading}
        <p>Loading config</p>
    {:else if error === undefined}
        <DirectorySelector directory={config["project_directory"]}
                           emitOn="config.select_project_directory"
                           listenOn="config.set_project_directory"
                           labelText="Project Directory"
                           helperText="Directory from which projects are discovered"
                           placeholder="~/Projects"
                           name="project_directory"
        />
        <Editors />
    {:else}
        <p>Something went wrong: {error}</p>
    {/if}
</div>
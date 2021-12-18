<script>
    import ProjectTile from "../components/ProjectTile.svelte";
    import {Headline} from "attractions";
    import {Accordion} from "svelte-collapsible";

    let error = undefined;
    let projects = undefined;
    let loading = true;

    // fetch all the projects. refresh is false because the initialization fetches projects with a full sync.
    window.backend.Projects.GetAll(false).then((data) => {
        console.log(data);
        projects = data;
        loading = false;
    }).catch((err) => {
        error = err;
        loading = false;
    })
</script>

<div>
    <Headline>Project List</Headline>
    {#if loading}
        <p>Loading</p>
    {:else if projects !== undefined && projects !== null}
        <Accordion>
            {#each projects as project}
                <ProjectTile project={project} projectDirectory="~/Projects"/>
            {/each}
        </Accordion>
    {:else if error === undefined}
        <p>No projects</p>
    {:else}
        <p>Something went wrong: {error}</p>
    {/if}
</div>

<style>
</style>

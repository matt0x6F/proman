<script>
    import { Label } from "attractions";
    import { Button } from "svelma";
    import {AccordionItem} from "svelte-collapsible";
    import {Icon} from "svelte-awesome";
    import {github} from "svelte-awesome/icons"

    // props
    export let project = undefined;
    export let projectDirectory = undefined;

    let hover = true;

    // opens a github project with the default browser
    const openGitHubProject = (event) => {
        event.preventDefault();
        // stops openProject from executing
        event.stopPropagation();

        // for some reason on:click seems to think the target is the underlying svg
        // this selects the parent element which is <a>
        window.wails.Events.Emit("OpenURL", event.target.parentElement.getAttribute("href"));
    }

    // opens a project with an IDE
    const openProject = (event) => {
        event.preventDefault();

        window.wails.Events.Emit("OpenProject", project.path);
    }

    const hashCode = (s) => {
        for(var i = 0, h = 0; i < s.length; i++)
            h = Math.imul(31, h) + s.charCodeAt(i) | 0;
        return h;
    }
</script>

<AccordionItem key={hashCode(projectDirectory + "/" + project.path)}>
    <div slot="header" class="project-tile-header">
        <Label class="project-tile-header-name">{project.path}</Label>
        <small class="project-tile-header-path">{projectDirectory}/{project.path}</small>
    </div>
    <div slot="body">
        {#if project.repository_urls !== undefined}
            <a href={project.repository_urls[0]} class="project-tile-github-link" on:click={openGitHubProject}>
                <Icon class="project-tile-github-logo" label="GitHub link" data={github} />
            </a>
        {:else}
            <small>No VCS providers detected</small>
        {/if}
    </div>
</AccordionItem>

<style>
    @use 'theme.css';

    div {
        margin-bottom: 1em;
    }

    small {
        font-size: .75em;
    }

    :global(.project-tile-header) {
        display: grid;
        grid-template-columns: [name] max-content [path] auto;
    }

    :global(.project-tile-header-name) {
        margin: auto .5em auto 0 !important;
    }

    :global(.project-tile-header-path) {
        margin: auto 0 auto 0;
        color: #888 !important;
        font-style: italic;
    }

    :global(.project-tile-github-link) {
        color: #000;
    }

    :global(.project-tile-github-link:hover) {
        color: #000;
    }
</style>
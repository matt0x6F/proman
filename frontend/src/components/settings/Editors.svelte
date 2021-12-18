<script>
    import {Button, ModalCard, Field, Input, Select, Icon} from 'svelma';
    import {Headline} from "attractions";

    let editors = undefined;
    let loading = true;
    let error = undefined;
    let selected = "0";

    let openModal = false;

    window.backend.EditorConfig.GetAll(false).then((data) => {
        console.log(data);
        editors = data;
        loading = false;
    }).catch((err) => {
        error = err;
        loading = false;
    })
</script>


<div>
    <Headline>Editor Settings</Headline>

    {#if loading}
        <p>Loading editor config</p>
    {:else if editors === null}
        <p>No editors configured</p>
    {:else if error === undefined}
        <Field>
            <Select placeholder="Select an editor">
                {#each Object.keys(editors) as key}
                    <option value={key}>{editors[key].name}</option>
                {/each}
            </Select>
        </Field>
    {:else}
        <p>Something went wrong: {error}</p>
    {/if}
    <Button class="add-editor" type="is-primary" size="is-small" on:click={() => openModal = true}>
        <Icon icon="plus" />
    </Button>
    <ModalCard bind:active={openModal} title="Add an Editor">
        <Field label="Editor Name" type={error === undefined ? null : "is-danger"} message={error === null ? helperText : error}>
            <Input placeholder="Name" />
        </Field>
        <Field label="Editor Path" type={error === undefined ? null : "is-danger"} message={error === null ? helperText : error}>
            <Input placeholder="Path" />
        </Field>
    </ModalCard>
</div>

<style>
    :global(.heading) {
        margin-top: .25em;
    }

    :global(.add-editor) {
        margin-top: 1em;
    }
</style>

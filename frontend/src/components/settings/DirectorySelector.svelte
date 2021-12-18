<script>
    import { Field, Input } from 'svelma';

    // the initial value
    export let directory = undefined;
    // the event to emit when the directory selector is clicked
    export let emitOn = undefined;
    // the event to listen for
    export let listenOn = undefined;

    export let labelText = undefined;
    export let helperText = undefined;
    export let placeholder = undefined;
    export let name = undefined;

    let error = null;

    // event handlers
    const selectProjectDir = () => {
        window.wails.Events.Emit(emitOn);
    }

    // the backend will update the value based on what the user chose
    window.wails.Events.On(listenOn, (dir, err) => {
        directory = dir;
        error = err;
    })
</script>

<Field label={labelText} type={error === null ? null : "is-danger"} message={error === null ? helperText : error}>
    <Input value={directory} on:click={selectProjectDir} placeholder={placeholder} />
</Field>
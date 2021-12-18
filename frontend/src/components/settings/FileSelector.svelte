<script>
    import {TextInput} from "carbon-components-svelte";

    // the initial value
    export let file = undefined;
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
    window.wails.Events.On(listenOn, (f, err) => {
        file = f;
        error = err;
    })
</script>

<TextInput
        labelText={labelText}
        helperText={helperText}
        placeholder={placeholder}
        value={file}
        name={name}
        on:click={selectProjectDir}
        invalid={error !== null}
        invalidText={error}
/>
<script>
	import 'bulma/css/bulma.css';

	import Header from './components/Header.svelte';

    // frontend router configuration
	// https://github.com/italypaleale/svelte-spa-router
	import Router from 'svelte-spa-router';
	// enable route-wrapping:
	// https://github.com/ItalyPaleAle/svelte-spa-router/blob/master/Advanced%20Usage.md#route-wrapping
	import {wrap} from 'svelte-spa-router/wrap'

	import Home from './views/Home.svelte';
	import Settings from "./views/Settings.svelte";

	const routes = {
		// Exact path
		'/': Home,
		// '/git/config': GitConfig,
		'/app/settings': Settings,
		// // catch all
		// '*': NotFound
	}

	let showSuccess = false;
	let title = '';
	let subtitle = '';
</script>

<main>
	<Header />
	<div>
		<Router {routes} />
	</div>
</main>

<style>
	main {
		display: grid;
		grid-template-rows: [nav] min-content [content] auto;
		grid-template-columns: [main] auto;
		font-family: "Ubuntu", sans-serif;
	}

	div {
		padding-left: .5em;
	}
</style>
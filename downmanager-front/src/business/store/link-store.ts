import { derived, writable, type Writable } from "svelte/store";
import type { Link } from "../model/link-model";

class LinkStore {
    constructor(
        private _links: Writable<Link[]> = writable([]),
        private _fetchingLinks: Writable<boolean> = writable(false)
    ) {}

    public get fetchingLinks() {
        return derived([this._fetchingLinks], ([$_fetchingLinks]) => $_fetchingLinks)
    }

    public get links() {
        return derived([this._links], ([$_links]) => $_links)
    }

    public async retrieveLinks() {
        this._fetchingLinks.set(true)
        try {
            const response = await fetch("http://localhost:8080/links", {
                method: 'GET'
            })
            const data = await response.json() as Link[]
            this._links.set(data)
            console.log(data)
        } catch (error) {
            console.error(error)
        } finally {
            this._fetchingLinks.set(false)
        }
    }
}

export const linkStore = new LinkStore()
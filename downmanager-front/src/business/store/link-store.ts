import { derived, writable, type Writable } from "svelte/store";
import type { Link } from "../model/link-model";
import { env } from "$env/dynamic/public";

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
            const response = await fetch(`${env.PUBLIC_BASE_API}/links`, {
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

    public async addLink(link: string) {
        const response = await fetch(`${env.PUBLIC_BASE_API}/links?link=${link}`, {
            method: 'POST'
        })
        if(response.ok) {
            const data = await response.json() as Link
            this._links.update(links => [...links, data])
        } else {
            console.log(`Error while creating link. Code : ${response.status} with message : ${response.statusText}`)
        }
    }

    public async removeLink(linkref: string) {
        const response = await fetch(`${env.PUBLIC_BASE_API}/links/${linkref}`, {
            method: 'DELETE'
        })
        if(response.ok) {
            this._links.update(links => links.filter(link => link.Ref != linkref))
        } else {
            console.log(`Error while creating link. Code : ${response.status} with message : ${response.statusText}`)
        }
    }
}

export const linkStore = new LinkStore()
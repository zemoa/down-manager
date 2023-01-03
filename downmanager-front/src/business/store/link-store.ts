import { derived, writable, type Writable } from "svelte/store";
import type { LinkItem } from "../model/link-model";
import { env } from "$env/dynamic/public";

class LinkStore {
    constructor(
        private _links: Writable<LinkItem[]> = writable([]),
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
            const data = await response.json() as LinkItem[]
            data.forEach(link => this.computePercent(link))
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
            const data = await response.json() as LinkItem
            this.computePercent(data)
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
            console.log(`Error while deleting link. Code : ${response.status} with message : ${response.statusText}`)
        }
    }

    public async startDownload(linkref: string) {
        this.changeDownloadState(linkref, true)
    }

    public async stopDownload(linkref: string) {
        this.changeDownloadState(linkref, false)
    }

    async changeDownloadState(linkref: string, start: boolean) {
        const action = (start? "start": "stop")
        const response = await fetch(`${env.PUBLIC_BASE_API}/links/${linkref}/${action}`, {
            method: 'PUT'
        })
        if(response.ok) {
            const data = await response.json() as LinkItem
            this.computePercent(data)
            this._links.update(links => {
                const link = links.find(link => linkref == link.Ref)
                if(link) {
                    link.Running = data.Running
                    return [...links]
                }
                return links
            })
        } else {
            console.log(`Error while calling action ${action} on link ${linkref}. Code : ${response.status} with message : ${response.statusText}`)
        }
    }

    private computePercent(link: LinkItem): LinkItem {
        link.Percent = Math.floor((link.Downloaded / link.Size) * 100)
        return link
    }
}

export const linkStore = new LinkStore()
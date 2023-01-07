import { derived, writable, type Writable } from "svelte/store";
import type { LinkItem, UpdateProgress } from "../model/link-model";
import { env } from "$env/dynamic/public";

class LinkStore {
    constructor(
        private _links: Writable<Map<string, LinkItem>> = writable(new Map()),
        private _fetchingLinks: Writable<boolean> = writable(false)
    ) {}

    public get fetchingLinks() {
        return derived([this._fetchingLinks], ([$_fetchingLinks]) => $_fetchingLinks)
    }

    public get links() {
        return derived([this._links], ([$_links]) => Array.from($_links.values()))
    }

    public async retrieveLinks() {
        this._fetchingLinks.set(true)
        try {
            const response = await fetch(`${env.PUBLIC_BASE_API}/links`, {
                method: 'GET'
            })
            const data = await response.json() as LinkItem[]
            data.forEach(link => this.computePercent(link))
            this._links.set(new Map(data.map(link => [link.Ref, link])))
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

            this._links.update(links => {
                links.set(data.Ref, data)
                return links
            })
        } else {
            console.log(`Error while creating link. Code : ${response.status} with message : ${response.statusText}`)
        }
    }

    public async removeLink(linkref: string) {
        const response = await fetch(`${env.PUBLIC_BASE_API}/links/${linkref}`, {
            method: 'DELETE'
        })
        if(response.ok) {
            this._links.update(links => {
                links.delete(linkref)
                return links
            })
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
                const link = links.get(linkref)
                if(link) {
                    link.Running = data.Running
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

    async updateProgress(linkRef: string, updateProgress: UpdateProgress) {
        this._links.update(links => {
            const link = links.get(linkRef)
            if(link) {
                link.Downloaded = updateProgress.Downloaded
                link.InError = updateProgress.InError
                link.Running = !updateProgress.Finished
                this.computePercent(link)
            }
            return links
        })
    }
}

export const linkStore = new LinkStore()
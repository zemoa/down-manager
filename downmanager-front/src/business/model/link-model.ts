export interface Link {
    Ref: string,
    Running: boolean,
    Link: string,
    InError: boolean,
    Size: number,
    Downloaded: number,
    Filename: string
}

export interface LinkItem extends Link {
    Percent: number
}

export interface UpdateProgress {
    InError: boolean,
    Downloaded: number,
    Finished: boolean,
}
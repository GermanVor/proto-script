type ValueOf<Obj> = Obj[keyof Obj];

type OneOnly<Obj, Key extends keyof Obj> = {
    [key in Exclude<keyof Obj, Key>]+?: undefined;
} & Pick<Obj, Key>;

type OneOfByKey<Obj> = {
    [key in keyof Obj]: OneOnly<Obj, key>;
};

export type OneOf<OneOfKey extends string, Obj> = {
    [key in OneOfKey]?: keyof Obj;
} & ValueOf<OneOfByKey<Obj>>;

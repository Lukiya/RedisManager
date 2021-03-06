import copy from 'copy-to-clipboard';

function createShiftArr(step: any) {

    var space = '    ';

    if (isNaN(parseInt(step))) {  // argument is string
        space = step;
    } else { // argument is integer
        switch (step) {
            case 1: space = ' '; break;
            case 2: space = '  '; break;
            case 3: space = '   '; break;
            case 4: space = '    '; break;
            case 5: space = '     '; break;
            case 6: space = '      '; break;
            case 7: space = '       '; break;
            case 8: space = '        '; break;
            case 9: space = '         '; break;
            case 10: space = '          '; break;
            case 11: space = '           '; break;
            case 12: space = '            '; break;
        }
    }

    var shift = ['\n']; // array of shifts
    for (var ix = 0; ix < 100; ix++) {
        shift.push(shift[ix] + space);
    }
    return shift;
}

const u = {
    STRING: "string",
    HASH: "hash",
    LIST: "list",
    SET: "set",
    ZSET: "zset",
    NONE: "none",
    CLIPBOARD_REDIS: "REDIS:",
    DefaultQuery: {
        serverID: '',
        nodeID: '',
        key: '',
        type: '',
        db: 0,
        cursor: 0,
        count: 50,
        keyword: '',
        all: false,
    },
    DefaultEntry: {
        Key: '',
        Type: '',
        Field: '',
        Index: 0,
        Score: 0,
        TTL: -1,
        Value: '',
    },
    LocalRootURL: () => {
        const r = process.env.NODE_ENV === "production" ? "/" : "http://localhost:16379/";
        return r;
    },
    IsXml: (str: string) => {
        const startPattern = /^\s*<[^>]+>/;
        const endPattern = /<\/[^>]+>\s*$/;
        return startPattern.test(str) && endPattern.test(str);
    },
    IsJson: (str: string) => {
        const pattern = /(^\s*\[[\s\S]*\]\s*$)|(^\s*\{[\s\S]*\}\s*$)/;
        return pattern.test(str);
    },
    FormatJson: (text: string) => {
        try {
            const step = '    ';

            if (typeof JSON === 'undefined') return text;

            if (typeof text === "string") return JSON.stringify(JSON.parse(text), null, step);
            if (typeof text === "object") return JSON.stringify(text, null, step);

            return text; // text is not string nor object
        } catch (err) {
            console.error(err);
            return text;
        }
    },
    MinifyJson: (text: string) => {
        try {
            if (typeof JSON === 'undefined') return text;

            return JSON.stringify(JSON.parse(text), null, 0);
        } catch (err) {
            console.error(err);
            return text;
        }
    },
    FormatXml: (text: string) => {
        try {
            const step = '    ';
            const shift1 = createShiftArr(step);

            var ar = text.replace(/>\s{0,}</g, "><")
                .replace(/</g, "~::~<")
                .replace(/\s*xmlns\:/g, "~::~xmlns:")
                .replace(/\s*xmlns\=/g, "~::~xmlns=")
                .split('~::~'),
                len = ar.length,
                inComment = false,
                deep = 0,
                str = '',
                ix = 0,
                shift = step ? createShiftArr(step) : shift1;

            for (ix = 0; ix < len; ix++) {
                // start comment or <![CDATA[...]]> or <!DOCTYPE //
                if (ar[ix].search(/<!/) > -1) {
                    str += shift[deep] + ar[ix];
                    inComment = true;
                    // end comment  or <![CDATA[...]]> //
                    if (ar[ix].search(/-->/) > -1 || ar[ix].search(/\]>/) > -1 || ar[ix].search(/!DOCTYPE/) > -1) {
                        inComment = false;
                    }
                } else
                    // end comment  or <![CDATA[...]]> //
                    if (ar[ix].search(/-->/) > -1 || ar[ix].search(/\]>/) > -1) {
                        str += ar[ix];
                        inComment = false;
                    } else
                        // <elm></elm> //
                        if (
                            /^<\w/.exec(ar[ix - 1])
                            &&
                            /^<\/\w/.exec(ar[ix])
                            // &&
                            // /^<[\w:\-\.\,]+/.exec(ar[ix - 1]) == /^<\/[\w:\-\.\,]+/.exec(ar[ix])[0].replace('/', '')
                        ) {
                            str += ar[ix];
                            if (!inComment) deep--;
                        } else
                            // <elm> //
                            if (ar[ix].search(/<\w/) > -1 && ar[ix].search(/<\//) == -1 && ar[ix].search(/\/>/) == -1) {
                                str = !inComment ? str += shift[deep++] + ar[ix] : str += ar[ix];
                            } else
                                // <elm>...</elm> //
                                if (ar[ix].search(/<\w/) > -1 && ar[ix].search(/<\//) > -1) {
                                    str = !inComment ? str += shift[deep] + ar[ix] : str += ar[ix];
                                } else
                                    // </elm> //
                                    if (ar[ix].search(/<\//) > -1) {
                                        str = !inComment ? str += shift[--deep] + ar[ix] : str += ar[ix];
                                    } else
                                        // <elm/> //
                                        if (ar[ix].search(/\/>/) > -1) {
                                            str = !inComment ? str += shift[deep] + ar[ix] : str += ar[ix];
                                        } else
                                            // <? xml ... ?> //
                                            if (ar[ix].search(/<\?/) > -1) {
                                                str += shift[deep] + ar[ix];
                                            } else
                                                // xmlns //
                                                if (ar[ix].search(/xmlns\:/) > -1 || ar[ix].search(/xmlns\=/) > -1) {
                                                    str += shift[deep] + ar[ix];
                                                }

                                                else {
                                                    str += ar[ix];
                                                }
            }

            return (str[0] == '\n') ? str.slice(1) : str;
        } catch (err) {
            console.error(err);
            return text;
        }
    },
    MinifyXml: (text: string, preserveComments: boolean) => {
        var str = preserveComments ? text
            : text.replace(/\<![ \r\n\t]*(--([^\-]|[\r\n]|-[^\-])*--[ \r\n\t]*)\>/g, "")
                .replace(/[ \r\n\t]{1,}xmlns/g, ' xmlns');
        return str.replace(/>\s{0,}</g, "><");
    },
    GetPageSize: () => {
        let r;
        const docHeight = document?.body?.clientHeight ?? 640;
        if (docHeight < 920) {
            r = 10;
        } else if (docHeight < 1210) {
            r = 20;
        } else {
            r = 30;
        }

        return r;
    },
    KeySorter: (a: any, b: any) => {
        const aType = typeof (a.Key);
        if (aType == "string") {
            return a.Key.localeCompare(b.Key);
        } else {
            return a.Key - b.Key;
        }
    },
    ValueSorter: (a: any, b: any) => {
        const aType = typeof (a.Value);
        if (aType == "string") {
            return a.Value.localeCompare(b.Value);
        } else {
            return a.Value - b.Value;
        }
    },
    OpenEditorForCreate: (params: any, key: any, type: string, dispatch: any) => {
        const payload = {
            ...params,
            entry: u.DefaultEntry,
            isNew: true,
            loading: false,
            keyEditorEnabled: true,
            valueEditorEnabled: true,
            fieldEditorEnabled: false,
            scoreEditorEnabled: false,
            indexEditorEnabled: false,
        };

        payload.entry.Key = key;

        switch (type) {
            case u.STRING:
                payload.entry.Type = u.STRING;
                break;
            case u.HASH:
                payload.fieldEditorEnabled = true;
                payload.entry.Type = u.HASH;
                break;
            case u.LIST:
                payload.indexEditorEnabled = true;
                payload.entry.Type = u.LIST;
                break;
            case u.SET:
                payload.entry.Type = u.SET;
                break;
            case u.ZSET:
                payload.scoreEditorEnabled = true;
                payload.entry.Type = u.ZSET;
                break;
        }

        dispatch({ type: "memberEditorVM/show", payload });
    },
    EntryToElement: (redisEntry: any) => {
        switch (redisEntry.Type) {
            case u.HASH:
                return {
                    Key: redisEntry.Field,
                    Value: redisEntry.Value,
                };
            case u.LIST:
                return {
                    Key: redisEntry.Index,
                    Value: redisEntry.Value,
                };
            case u.SET:
                return {
                    Key: redisEntry.Value,
                    // Value: redisEntry.Value,
                };
            case u.ZSET:
                return {
                    Key: redisEntry.Value,
                    Value: redisEntry.Score,
                };
        }
    },
    CopyToClipboard: (content: any) => {
        copy(u.CLIPBOARD_REDIS + content);
    },
    Base64ToBytes: (base64: any) => {
        const binStr = window.atob(base64);
        const len = binStr.length;
        let bytes = new Array(len);
        for (var i = 0; i < len; i++) {
            bytes[i] = binStr.charCodeAt(i);
        }
        return bytes;
    },
    IsPresent: (input: any) => {
        const a = input !== undefined && input !== null;
        if (a && typeof (input) === "string") {
            return input.trim().length > 0;
        } else {
            return a;
        }
    },
    IsMissing: (input: any) => {
        const a = input === undefined || input === null;
        if (!a && typeof (input) === "string") {
            return input.trim().length == 0;
        } else {
            return a;
        }
    },
}
export default u
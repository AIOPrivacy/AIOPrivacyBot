const webUrl = window.location.href;
const headline = document.title;
const host = location.host;

const sites = [
    { "host": "blog.csdn.net", "el": "article.baidu_pl", "cut_str": "_" },
    { "host": "www.jianshu.com", "el": "article._2rhmJa", "cut_str": " - " },
    { "host": "juejin.cn", "el": ".article-viewer.markdown-body.result", "cut_str": " - " },
    { "host": "zhuanlan.zhihu.com", "el": ".Post-RichTextContainer", "cut_str": " - " },
    { "host": "www.cnblogs.com", "el": "#cnblogs_post_body", "cut_str": " - " },
    { "host": "www.jb51.net", "el": "#content", "cut_str": "_" },
    { "host": "blog.51cto.com", "el": "#result", "cut_str": "_" },
    { "host": "www.pianshen.com", "el": ".blogpost-body", "cut_str": " - " },
    { "host": "www.360doc.com", "el": "#artContent", "cut_str": "" },
    { "host": "baijiahao.baidu.com", "el": "div[data-testid='article']", "cut_str": "" },
    { "host": "jingyan.baidu.com", "el": ".exp-content-outer", "cut_str": "-" },
    { "host": "www.52pojie.cn", "el": ".t_f", "cut_str": " - " },
    { "host": "cloud.tencent.com", "el": ".mod-content__markdown", "cut_str": "-" },
    { "host": "developer.aliyun.com", "el": ".content-wrapper", "cut_str": "-" },
    { "host": "huaweicloud.csdn.net", "el": ".main-content", "cut_str": "_" },
    { "host": "www.bilibili.com", "el": "#read-article-holder", "cut_str": " - " },
    { "host": "weibo.com", "el": ".main_editor", "cut_str": "" },
    { "host": "www.weibo.com", "el": ".main_editor", "cut_str": "" },
    { "host": "mp.weixin.qq.com", "el": "#js_content", "cut_str": "" },
    { "host": "segmentfault.com", "el": ".article.fmt.article-content", "cut_str": "- SegmentFault 思否" },
    { "host": "www.qinglite.cn", "el": ".markdown-body", "cut_str": "-" },
    { "host": "www.manongjc.com", "el": "#code_example", "cut_str": " - " }

]

const cutTitle = (title, cut_str) => {
    try {
        const newTitle = title.split(cut_str)[0];
        return newTitle;
    }
    catch (e) {
        console.log(e);
        return title;
    }
}

const domToNode = (domNode) => {
    if (domNode.nodeType == domNode.TEXT_NODE) {
        return domNode.data;
    }
    if (domNode.nodeType != domNode.ELEMENT_NODE) {
        return false;
    }
    let nodeElement = {};
    nodeElement.tag = domNode.tagName.toLowerCase();
    for (const attr of domNode.attributes) {
        if (attr.name == 'href' || attr.name == 'src') {
            if (!nodeElement.attrs) {
                nodeElement.attrs = {};
            }
            nodeElement.attrs[attr.name] = attr.value;
        }
    }
    if (domNode.childNodes.length > 0) {
        nodeElement.children = [];
        for (const child of domNode.childNodes) {
            nodeElement.children.push(domToNode(child));
        }
    }
    return nodeElement;
}

const getData = () => {
    let new_headline;

    for (const site of sites) {
        if (!host.endsWith(site.host)) continue;
        const cut = site.cut_str;

        if (cut != '') {
            new_headline = cutTitle(headline, cut);
        } else {
            new_headline = document.title;
        }

        const ele = document.querySelector(site.el)
        const rootNode = domToNode(ele);
        rootNode.children.push(`\n\n本文转自 ${webUrl} ，如有侵权，请联系删除。`)
        rootNode.children.push(`\n本文接到 @AIOPrivacyBot 用户的Inline请求，用户要求保护其隐私，因此将其转换为Telegraph文档发送。`)
        const data = {
            title: new_headline,
            node: JSON.stringify(rootNode.children),
        };

        return data;
    }

    return null;
}

getData();

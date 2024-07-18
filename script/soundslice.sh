#!/bin/bash

source .env

gp_file=${1}

enable_embed() {
    local score_id=${1}
    curl -XPOST -sSL "https://www.soundslice.com/api/v1/scores/${score_id}/" \
        -H "referer: https://www.soundslice.com" \
        -H "cookie: sesn=${SESN}" \
        -H "content-type: application/x-www-form-urlencoded" \
        -d 'embed_status=4'
}

disable_embed() {
    local score_id=${1}
    curl -XPOST -sSL "https://www.soundslice.com/api/v1/scores/${score_id}/" \
        -H "referer: https://www.soundslice.com" \
        -H "cookie: sesn=${SESN}" \
        -H "content-type: application/x-www-form-urlencoded" \
        -d 'embed_status=1'
}

get_score_id() {
    local slice_id=${1}
    curl -XGET -sSL https://www.soundslice.com/slices/${slice_id}/edit/scoredata/ \
        -H "referer: https://www.soundslice.com" \
        -H "cookie: sesn=${SESN}" \
        | jq -r '.slug'
}

create_notation() {
    curl -XPOST -sSL https://www.soundslice.com/manage/create-via-import/ \
        -H "referer: https://www.soundslice.com" \
        -H "cookie: sesn=${SESN}" | \
        grep title-practice-lists | \
        sed 's/.*slice="\(.*\)".*/\1/'
}

upload_notation() {
    local slice_id=${1}
    local file=${2}
    curl -XPOST -sSL "https://www.soundslice.com/api/v1/slices/${slice_id}/notation/" \
        -H "content-type: multipart/form-data" \
        -H "referer: https://www.soundslice.com" \
        -H "cookie: sesn=${SESN}" \
        -F "type=application/octet-stream" \
        -F "score=@${file}" | jq -r '.name'
}

delete_notation() {
    local slice_id=${1}
    curl -XPOST -sSL "https://www.soundslice.com/api/v1/slices/delete-multiple/" \
        -H "referer: https://www.soundslice.com" \
        -H "cookie: sesn=${SESN}" \
        -H 'content-type: application/x-www-form-urlencoded' \
        -d "ids=${slice_id}"
}

list_scores() {
    curl -XGET -sSL "https://www.soundslice.com/" \
        -H "referer: https://www.soundslice.com" \
        -H "cookie: sesn=${SESN}"
}

generate_html() {
    local slice_id=${1}
    cat <<EOF > index.html
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
    <iframe id="soundsliceiframe" src="https://www.soundslice.com/slices/${slice_id}/embed/" width="100%" height="400"></iframe>
</body>
<script type="text/javascript" src="/js/autoheight.min.js"></script>
<script type="text/javascript">
    autoheight('soundsliceiframe');
</script>
</html>
EOF
}

# slice_id=$(cat .cache_slice_id)
# disable_embed $(get_score_id ${slice_id})
# delete_notation ${slice_id}

# slice_id=$(create_notation)
# upload_notation ${slice_id} "${gp_file}"
# enable_embed $(get_score_id ${slice_id})
# generate_html ${slice_id}

# echo ${slice_id} > .cache_slice_id

list_scores

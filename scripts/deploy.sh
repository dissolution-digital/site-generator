BUCKET='SOMETHING'
CLOUDFRONT='SOMETHING'
echo $BUCKET

aws s3 cp . $BUCKET --recursive

aws cloudfront create-invalidation --distribution-id $CLOUDFRONT --paths "/*"

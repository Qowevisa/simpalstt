### Go to localhost:8080

```graphql
query {
  data(title:"se vinde", limit:5) {
    data {
      id
      title {
        ru
        ro
      }
      categories {
        subcategory
      }
      type
      posted
    }
    nextPageToken
  }
}
```


```graphql
query {
  data(title:"se vinde", limit:5, pageToken:"1488276296.990566") {
    data {
      id
      title {
        ru
        ro
      }
      categories {
        subcategory
      }
      type
      posted
    }
    nextPageToken
  }
}
```

```graphql
query {
  aggregateSubcategory(subcategory:"") {
    category {
      subcategory
    }
    count
  }
}
```

```graphql
query {
  aggregateSubcategory(subcategory:"1404") {
    category {
      subcategory
    }
    count
  }
}
```

```graphql
query {
  aggregateSubcategory(subcategory:"1300") {
    category {
      subcategory
    }
    count
  }
}
```

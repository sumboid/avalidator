import {gql, useQuery} from '@apollo/client';
import React, {memo} from 'react';

const Page: React.FC = () => {
  const {loading, error, data} = useQuery(gql`
    query Users {
      users {
        name
        group
        email
      }
    }
  `);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error! {error.message}</div>;
  }

  console.log(data);

  return <div>Success!</div>;
};

export default memo(Page);

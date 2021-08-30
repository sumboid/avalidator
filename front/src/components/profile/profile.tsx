import React, {memo} from 'react';
import {Avatar, Typography, makeStyles} from '@material-ui/core';

const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
    width: '100%',
  },
  avatar: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    padding: '10px 0',
  },
  name: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    padding: '5px 0',
  },
  email: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  },
  large: {
    width: theme.spacing(7),
    height: theme.spacing(7),
  },
}));

type Props = {
  name: string;
  email: string;
};

export default memo(({name, email}: Props) => {
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <div className={classes.avatar}>
        <Avatar
          src={`https://robohash.org/${email}`}
          className={classes.large}
        />
      </div>
      <div className={classes.name}>
        <Typography variant="h6">{name}</Typography>
      </div>
      <div className={classes.email}>{email}</div>
    </div>
  );
});

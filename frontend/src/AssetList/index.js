/**
 * Copyright 2019 Atelier Disko. All rights reserved. This source
 * code is distributed under the terms of the BSD 3-Clause License.
 */

import React from 'react';
import './AssetList.css';
import Meta from '../Meta';
import Heading from '../Heading';
import filesize from 'filesize';

function AssetList(props) {
  const imageFileTypes = ['png', 'jpg', 'jpeg'];

  return (
    <>
      {props.assets &&
        props.assets.map(a => {
          return (
            <div className="asset-list__asset" key={a.name}>
              <Heading level="beta" isJumptarget={true} docTitle="Assets">
                {a.name}
              </Heading>

              {imageFileTypes.some(v => {
                return a.url.indexOf(v) >= 0;
              }) && (
                <img className="asset-list__asset-image" src={`/api/v1/tree/${a.url}?v=${props.source}`} alt={a.name} />
              )}

              <div className="asset-list__asset-meta">
                <Meta title="Download">
                  <a className="asset-list__asset-download" href={`/api/v1/tree/${a.url}?v=${props.source}`} download>
                    {a.name}
                  </a>
                </Meta>
                <Meta title="Size">{filesize(a.size)}</Meta>
                <Meta title="Last Modified">{new Date(a.modified * 1000).toLocaleString()}</Meta>
              </div>
            </div>
          );
        })}
    </>
  );
}

export default AssetList;

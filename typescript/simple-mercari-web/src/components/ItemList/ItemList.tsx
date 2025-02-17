import React, { useEffect, useState } from 'react';

interface Item {
  id: number;
  name: string;
  category: string;
  image_filename: string;
};

const server = process.env.API_URL || 'http://127.0.0.1:9000';

interface Prop {
  reload?: boolean;
  onLoadCompleted?: () => void;
}

export const ItemList: React.FC<Prop> = (props) => {
  const { reload = true, onLoadCompleted } = props;
  const [items, setItems] = useState<Item[]>([])
  const fetchItems = () => {
    fetch(server.concat('/items'),
      {
        method: 'GET',
        mode: 'cors',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        },
      })
      .then(response => response.json())
      .then(data => {
        console.log('GET success:', data);
        setItems(data.items);
        onLoadCompleted && onLoadCompleted();
      })
      .catch(error => {
        console.error('GET error:', error)
      })
  }

  useEffect(() => {
    if (reload) {
      fetchItems();
    }
  }, [reload]);

  return (
    <div>
      {items.map((item) => {
        return (
          <div key={item.id} className='ItemList'>
            <div className='ItemImage'>
              <img src={`http://127.0.0.1:9000/image/${item.id}.jpg`} alt=""/>
            </div>
            <div className='ItemName'>
              <span><b>Name: {item.name}</b></span>
            </div>
              <br/>
            <p>
              <div className='ItemInfo'>
                <span><b>Category: {item.category}</b></span><br/>
              </div>
            </p>
          </div>
        )
      })}
    </div>
  )
};
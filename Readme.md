## Example Create Meeting JSON Data [input]

    {
    	"title": "Cool2",
    	"participants": [
    		{"name": "Sharath", "email": "tnfssc@gmail.com", "rsvp": "Yes"},
    		{"name": "Chandra", "email": "tnfssc1@gmail.com", "rsvp": "Yes"},
    		{"name": "Nikhil", "email": "tnfssc2@gmail.com", "rsvp": "No"},
    		{"name": "Geeks", "email": "tnfssc3@gmail.com", "rsvp": "Yes"}
    	],
    	"startTime": "2020-06-19T14:01:42.240Z",
    	"endTime": "2020-06-20T14:01:42.240Z"
    }

## Example Response Meeting JSON Data [output]

    {
        "_id": "5f6bfd8d5c1590406b255bae",
        "title": "Cool2",
        "participants": [
            {
                "name": "Sharath",
                "email": "tnfssc@gmail.com",
                "rsvp": "Yes"
            },
            {
                "name": "Chandra",
                "email": "tnfssc1@gmail.com",
                "rsvp": "Yes"
            },
            {
                "name": "Nikhil",
                "email": "tnfssc2@gmail.com",
                "rsvp": "No"
            },
            {
                "name": "Geeks",
                "email": "tnfssc3@gmail.com",
                "rsvp": "Yes"
            }
        ],
        "startTime": {
            "T": 1592575302,
            "I": 0
        },
        "endTime": {
            "T": 1592661702,
            "I": 0
        },
        "createdAt": {
            "T": 1600912781,
            "I": 0
        }
    }

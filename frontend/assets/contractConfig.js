// contractConfig.js
export const contractAddress = "0xfEbc1138f0b13F35BE8F7E820A6f486368E511CC";

export const contractABI = [
  {
    inputs: [
      {
        internalType: "string",
        name: "_roomId",
        type: "string",
      },
      {
        internalType: "string",
        name: "_sessionId",
        type: "string",
      },
      {
        internalType: "string",
        name: "_ingestUrl",
        type: "string",
      },
      {
        internalType: "string",
        name: "_streamKey",
        type: "string",
      },
      {
        internalType: "string",
        name: "_host",
        type: "string",
      },
      {
        internalType: "string",
        name: "_title",
        type: "string",
      },
    ],
    name: "createLivestream",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "string",
        name: "_roomId",
        type: "string",
      },
    ],
    name: "endLivestream",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: false,
        internalType: "string",
        name: "roomId",
        type: "string",
      },
      {
        indexed: false,
        internalType: "string",
        name: "sessionId",
        type: "string",
      },
      {
        indexed: false,
        internalType: "string",
        name: "ingestUrl",
        type: "string",
      },
      {
        indexed: false,
        internalType: "string",
        name: "streamKey",
        type: "string",
      },
      {
        indexed: false,
        internalType: "string",
        name: "title",
        type: "string",
      },
      {
        indexed: false,
        internalType: "string",
        name: "host",
        type: "string",
      },
    ],
    name: "LivestreamCreated",
    type: "event",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: false,
        internalType: "string",
        name: "roomId",
        type: "string",
      },
      {
        indexed: false,
        internalType: "string",
        name: "sessionId",
        type: "string",
      },
    ],
    name: "LivestreamEnded",
    type: "event",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: false,
        internalType: "string",
        name: "roomId",
        type: "string",
      },
      {
        indexed: false,
        internalType: "string",
        name: "newTitle",
        type: "string",
      },
    ],
    name: "TitleUpdated",
    type: "event",
  },
  {
    inputs: [
      {
        internalType: "string",
        name: "_roomId",
        type: "string",
      },
      {
        internalType: "string",
        name: "_newTitle",
        type: "string",
      },
    ],
    name: "updateTitle",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "string",
        name: "_roomId",
        type: "string",
      },
    ],
    name: "getLivestreamInfo",
    outputs: [
      {
        internalType: "string",
        name: "sessionId",
        type: "string",
      },
      {
        internalType: "string",
        name: "ingestUrl",
        type: "string",
      },
      {
        internalType: "string",
        name: "streamKey",
        type: "string",
      },
      {
        internalType: "string",
        name: "host",
        type: "string",
      },
      {
        internalType: "string",
        name: "title",
        type: "string",
      },
      {
        internalType: "string",
        name: "status",
        type: "string",
      },
      {
        internalType: "uint256",
        name: "createdAt",
        type: "uint256",
      },
    ],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "string",
        name: "",
        type: "string",
      },
    ],
    name: "livestreams",
    outputs: [
      {
        internalType: "string",
        name: "roomId",
        type: "string",
      },
      {
        internalType: "string",
        name: "sessionId",
        type: "string",
      },
      {
        internalType: "string",
        name: "ingestUrl",
        type: "string",
      },
      {
        internalType: "string",
        name: "streamKey",
        type: "string",
      },
      {
        internalType: "string",
        name: "host",
        type: "string",
      },
      {
        internalType: "string",
        name: "title",
        type: "string",
      },
      {
        internalType: "string",
        name: "status",
        type: "string",
      },
      {
        internalType: "uint256",
        name: "createdAt",
        type: "uint256",
      },
    ],
    stateMutability: "view",
    type: "function",
  },
];
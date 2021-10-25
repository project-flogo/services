package task

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/project-flogo/flow/state"
)

var stepData1 = `[
  {
    "id": 0,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": true,
        "flowURI": "res://flow:flow3",
        "subflowId": 0,
        "taskId": "SharedData",
        "status": 100,
        "attrs": null,
        "tasks": {
          "SharedData": {
            "change": 1,
            "status": 20,
            "input": null
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "1": {
        "change": 0,
        "subflowId": 0,
        "taskId": "SharedData"
      }
    }
  },
  {
    "id": 1,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "SharedData",
        "status": 0,
        "attrs": {
          "TIB_Flow:123": "123"
        },
        "tasks": {
          "SharedData": {
            "change": 2,
            "status": 40,
            "input": {
              "input": {
                "data": "123"
              },
              "key": "123"
            }
          },
          "SharedData1": {
            "change": 1,
            "status": 20,
            "input": null
          }
        },
        "links": {
          "0": {
            "change": 1,
            "status": 2,
            "from": "SharedData",
            "to": "SharedData1"
          }
        },
        "returnData": null
      }
    },
    "queueChanges": {
      "1": {
        "change": 2,
        "subflowId": 0,
        "taskId": "SharedData"
      },
      "2": {
        "change": 0,
        "subflowId": 0,
        "taskId": "SharedData1"
      }
    }
  },
  {
    "id": 2,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "SharedData1",
        "status": 0,
        "attrs": {
          "_A.SharedData1.output": {
            "data": "123",
            "exist": true
          }
        },
        "tasks": {
          "SharedData1": {
            "change": 1,
            "status": 20,
            "input": {
              "key": "123"
            }
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "2": {
        "change": 2,
        "subflowId": 0,
        "taskId": "SharedData1"
      },
      "3": {
        "change": 0,
        "subflowId": 0,
        "taskId": "SharedData1"
      }
    }
  },
  {
    "id": 3,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "SharedData1",
        "status": 0,
        "attrs": {
          "_A.SharedData1.output": {
            "data": "123",
            "exist": true
          }
        },
        "tasks": {
          "SharedData1": {
            "change": 1,
            "status": 20,
            "input": {
              "key": "123"
            }
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "3": {
        "change": 2,
        "subflowId": 0,
        "taskId": "SharedData1"
      },
      "4": {
        "change": 0,
        "subflowId": 0,
        "taskId": "SharedData1"
      }
    }
  },
  {
    "id": 4,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "SharedData1",
        "status": 0,
        "attrs": {
          "_A.SharedData1.output": {
            "data": "123",
            "exist": true
          }
        },
        "tasks": {
          "SharedData1": {
            "change": 1,
            "status": 20,
            "input": {
              "key": "123"
            }
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "4": {
        "change": 2,
        "subflowId": 0,
        "taskId": "SharedData1"
      },
      "5": {
        "change": 0,
        "subflowId": 0,
        "taskId": "SharedData1"
      }
    }
  },
  {
    "id": 5,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "SharedData1",
        "status": 0,
        "attrs": {
          "_A.SharedData1.output": {
            "data": "123",
            "exist": true
          }
        },
        "tasks": {
          "SharedData1": {
            "change": 1,
            "status": 20,
            "input": {
              "key": "123"
            }
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "5": {
        "change": 2,
        "subflowId": 0,
        "taskId": "SharedData1"
      },
      "6": {
        "change": 0,
        "subflowId": 0,
        "taskId": "SharedData1"
      }
    }
  },
  {
    "id": 6,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "SharedData1",
        "status": 0,
        "attrs": {
          "_A.SharedData1.output": {
            "data": "123",
            "exist": true
          }
        },
        "tasks": {
          "SharedData1": {
            "change": 1,
            "status": 20,
            "input": {
              "key": "123"
            }
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "6": {
        "change": 2,
        "subflowId": 0,
        "taskId": "SharedData1"
      },
      "7": {
        "change": 0,
        "subflowId": 0,
        "taskId": "SharedData1"
      }
    }
  },
  {
    "id": 7,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "SharedData1",
        "status": 0,
        "attrs": null,
        "tasks": {
          "Accumulate": {
            "change": 1,
            "status": 20,
            "input": null
          },
          "SharedData1": {
            "change": 2,
            "status": 40,
            "input": {
              "key": "123"
            }
          }
        },
        "links": {
          "0": {
            "change": 2,
            "status": 0,
            "from": "",
            "to": ""
          },
          "1": {
            "change": 1,
            "status": 2,
            "from": "SharedData1",
            "to": "Accumulate"
          }
        },
        "returnData": null
      }
    },
    "queueChanges": {
      "7": {
        "change": 2,
        "subflowId": 0,
        "taskId": "SharedData1"
      },
      "8": {
        "change": 0,
        "subflowId": 0,
        "taskId": "Accumulate"
      }
    }
  },
  {
    "id": 8,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "Accumulate",
        "status": 0,
        "attrs": {
          "_A.Accumulate": [
            {
              "output": {
                "data": "123",
                "exist": true
              }
            }
          ]
        },
        "tasks": {
          "Accumulate": {
            "change": 1,
            "status": 20,
            "input": {
              "key": "123"
            }
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "8": {
        "change": 2,
        "subflowId": 0,
        "taskId": "Accumulate"
      },
      "9": {
        "change": 0,
        "subflowId": 0,
        "taskId": "Accumulate"
      }
    }
  },
  {
    "id": 9,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "Accumulate",
        "status": 0,
        "attrs": {
          "_A.Accumulate": [
            {
              "output": {
                "data": "123",
                "exist": true
              }
            },
            {
              "output": {
                "data": "123",
                "exist": true
              }
            }
          ]
        },
        "tasks": {
          "Accumulate": {
            "change": 1,
            "status": 20,
            "input": {
              "key": "123"
            }
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "10": {
        "change": 0,
        "subflowId": 0,
        "taskId": "Accumulate"
      },
      "9": {
        "change": 2,
        "subflowId": 0,
        "taskId": "Accumulate"
      }
    }
  },
  {
    "id": 10,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "Accumulate",
        "status": 0,
        "attrs": {
          "_A.Accumulate": [
            {
              "output": {
                "data": "123",
                "exist": true
              }
            },
            {
              "output": {
                "data": "123",
                "exist": true
              }
            },
            {
              "output": {
                "data": "123",
                "exist": true
              }
            }
          ]
        },
        "tasks": {
          "Accumulate": {
            "change": 1,
            "status": 20,
            "input": {
              "key": "123"
            }
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "10": {
        "change": 2,
        "subflowId": 0,
        "taskId": "Accumulate"
      },
      "11": {
        "change": 0,
        "subflowId": 0,
        "taskId": "Accumulate"
      }
    }
  },
  {
    "id": 11,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "Accumulate",
        "status": 0,
        "attrs": {
          "_A.Accumulate": [
            {
              "output": {
                "data": "123",
                "exist": true
              }
            },
            {
              "output": {
                "data": "123",
                "exist": true
              }
            },
            {
              "output": {
                "data": "123",
                "exist": true
              }
            },
            {
              "output": {
                "data": "123",
                "exist": true
              }
            }
          ]
        },
        "tasks": {
          "Accumulate": {
            "change": 1,
            "status": 20,
            "input": {
              "key": "123"
            }
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "11": {
        "change": 2,
        "subflowId": 0,
        "taskId": "Accumulate"
      },
      "12": {
        "change": 0,
        "subflowId": 0,
        "taskId": "Accumulate"
      }
    }
  },
  {
    "id": 12,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "Accumulate",
        "status": 0,
        "attrs": {
          "_A.Accumulate": [
            {
              "output": {
                "data": "123",
                "exist": true
              }
            },
            {
              "output": {
                "data": "123",
                "exist": true
              }
            },
            {
              "output": {
                "data": "123",
                "exist": true
              }
            },
            {
              "output": {
                "data": "123",
                "exist": true
              }
            },
            {
              "output": {
                "data": "123",
                "exist": true
              }
            }
          ]
        },
        "tasks": {
          "Accumulate": {
            "change": 1,
            "status": 20,
            "input": {
              "key": "123"
            }
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "12": {
        "change": 2,
        "subflowId": 0,
        "taskId": "Accumulate"
      },
      "13": {
        "change": 0,
        "subflowId": 0,
        "taskId": "Accumulate"
      }
    }
  },
  {
    "id": 13,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "Accumulate",
        "status": 0,
        "attrs": null,
        "tasks": {
          "Accumulate": {
            "change": 2,
            "status": 40,
            "input": {
              "key": "123"
            }
          },
          "Dowhile": {
            "change": 1,
            "status": 20,
            "input": null
          }
        },
        "links": {
          "1": {
            "change": 2,
            "status": 0,
            "from": "",
            "to": ""
          },
          "2": {
            "change": 1,
            "status": 2,
            "from": "Accumulate",
            "to": "Dowhile"
          }
        },
        "returnData": null
      }
    },
    "queueChanges": {
      "13": {
        "change": 2,
        "subflowId": 0,
        "taskId": "Accumulate"
      },
      "14": {
        "change": 0,
        "subflowId": 0,
        "taskId": "Dowhile"
      }
    }
  },
  {
    "id": 14,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "Dowhile",
        "status": 0,
        "attrs": {
          "_A.Dowhile.output": {
            "data": "123",
            "exist": true
          }
        },
        "tasks": {
          "Dowhile": {
            "change": 1,
            "status": 20,
            "input": {
              "key": "123"
            }
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "14": {
        "change": 2,
        "subflowId": 0,
        "taskId": "Dowhile"
      },
      "15": {
        "change": 0,
        "subflowId": 0,
        "taskId": "Dowhile"
      }
    }
  },
  {
    "id": 15,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "Dowhile",
        "status": 0,
        "attrs": {
          "_A.Dowhile.output": {
            "data": "123",
            "exist": true
          }
        },
        "tasks": {
          "Dowhile": {
            "change": 1,
            "status": 20,
            "input": {
              "key": "123"
            }
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "15": {
        "change": 2,
        "subflowId": 0,
        "taskId": "Dowhile"
      },
      "16": {
        "change": 0,
        "subflowId": 0,
        "taskId": "Dowhile"
      }
    }
  },
  {
    "id": 16,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "Dowhile",
        "status": 0,
        "attrs": {
          "_A.Dowhile.output": {
            "data": "123",
            "exist": true
          }
        },
        "tasks": {
          "Dowhile": {
            "change": 1,
            "status": 20,
            "input": {
              "key": "123"
            }
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "16": {
        "change": 2,
        "subflowId": 0,
        "taskId": "Dowhile"
      },
      "17": {
        "change": 0,
        "subflowId": 0,
        "taskId": "Dowhile"
      }
    }
  },
  {
    "id": 17,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "Dowhile",
        "status": 0,
        "attrs": {
          "_A.Dowhile.output": {
            "data": "123",
            "exist": true
          }
        },
        "tasks": {
          "Dowhile": {
            "change": 1,
            "status": 20,
            "input": {
              "key": "123"
            }
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "17": {
        "change": 2,
        "subflowId": 0,
        "taskId": "Dowhile"
      },
      "18": {
        "change": 0,
        "subflowId": 0,
        "taskId": "Dowhile"
      }
    }
  },
  {
    "id": 18,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "Dowhile",
        "status": 0,
        "attrs": {
          "_A.Dowhile.output": {
            "data": "123",
            "exist": true
          }
        },
        "tasks": {
          "Dowhile": {
            "change": 1,
            "status": 20,
            "input": {
              "key": "123"
            }
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "18": {
        "change": 2,
        "subflowId": 0,
        "taskId": "Dowhile"
      },
      "19": {
        "change": 0,
        "subflowId": 0,
        "taskId": "Dowhile"
      }
    }
  },
  {
    "id": 19,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "Dowhile",
        "status": 0,
        "attrs": {
          "_A.Dowhile.output": {
            "data": "123",
            "exist": true
          }
        },
        "tasks": {
          "CopyOfDowhile": {
            "change": 1,
            "status": 20,
            "input": null
          },
          "Dowhile": {
            "change": 2,
            "status": 40,
            "input": {
              "key": "123"
            }
          }
        },
        "links": {
          "2": {
            "change": 2,
            "status": 0,
            "from": "",
            "to": ""
          },
          "3": {
            "change": 1,
            "status": 2,
            "from": "Dowhile",
            "to": "CopyOfDowhile"
          }
        },
        "returnData": null
      }
    },
    "queueChanges": {
      "19": {
        "change": 2,
        "subflowId": 0,
        "taskId": "Dowhile"
      },
      "20": {
        "change": 0,
        "subflowId": 0,
        "taskId": "CopyOfDowhile"
      }
    }
  },
  {
    "id": 20,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "CopyOfDowhile",
        "status": 0,
        "attrs": {
          "_A.CopyOfDowhile": [
            {
              "output": {
                "data": "123",
                "exist": true
              }
            }
          ],
          "_A.CopyOfDowhile.output": {
            "data": "123",
            "exist": true
          }
        },
        "tasks": {
          "CopyOfDowhile": {
            "change": 1,
            "status": 20,
            "input": {
              "key": "123"
            }
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "20": {
        "change": 2,
        "subflowId": 0,
        "taskId": "CopyOfDowhile"
      },
      "21": {
        "change": 0,
        "subflowId": 0,
        "taskId": "CopyOfDowhile"
      }
    }
  },
  {
    "id": 21,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "CopyOfDowhile",
        "status": 0,
        "attrs": {
          "_A.CopyOfDowhile": [
            {
              "output": {
                "data": "123",
                "exist": true
              }
            },
            {
              "output": {
                "data": "123",
                "exist": true
              }
            }
          ],
          "_A.CopyOfDowhile.output": {
            "data": "123",
            "exist": true
          }
        },
        "tasks": {
          "CopyOfDowhile": {
            "change": 1,
            "status": 20,
            "input": {
              "key": "123"
            }
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "21": {
        "change": 2,
        "subflowId": 0,
        "taskId": "CopyOfDowhile"
      },
      "22": {
        "change": 0,
        "subflowId": 0,
        "taskId": "CopyOfDowhile"
      }
    }
  },
  {
    "id": 22,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "CopyOfDowhile",
        "status": 0,
        "attrs": {
          "_A.CopyOfDowhile": [
            {
              "output": {
                "data": "123",
                "exist": true
              }
            },
            {
              "output": {
                "data": "123",
                "exist": true
              }
            },
            {
              "output": {
                "data": "123",
                "exist": true
              }
            }
          ],
          "_A.CopyOfDowhile.output": {
            "data": "123",
            "exist": true
          }
        },
        "tasks": {
          "CopyOfDowhile": {
            "change": 1,
            "status": 20,
            "input": {
              "key": "123"
            }
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "22": {
        "change": 2,
        "subflowId": 0,
        "taskId": "CopyOfDowhile"
      },
      "23": {
        "change": 0,
        "subflowId": 0,
        "taskId": "CopyOfDowhile"
      }
    }
  },
  {
    "id": 23,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "CopyOfDowhile",
        "status": 0,
        "attrs": {
          "_A.CopyOfDowhile": [
            {
              "output": {
                "data": "123",
                "exist": true
              }
            },
            {
              "output": {
                "data": "123",
                "exist": true
              }
            },
            {
              "output": {
                "data": "123",
                "exist": true
              }
            },
            {
              "output": {
                "data": "123",
                "exist": true
              }
            }
          ],
          "_A.CopyOfDowhile.output": {
            "data": "123",
            "exist": true
          }
        },
        "tasks": {
          "CopyOfDowhile": {
            "change": 1,
            "status": 20,
            "input": {
              "key": "123"
            }
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "23": {
        "change": 2,
        "subflowId": 0,
        "taskId": "CopyOfDowhile"
      },
      "24": {
        "change": 0,
        "subflowId": 0,
        "taskId": "CopyOfDowhile"
      }
    }
  },
  {
    "id": 24,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "CopyOfDowhile",
        "status": 0,
        "attrs": {
          "_A.CopyOfDowhile": [
            {
              "output": {
                "data": "123",
                "exist": true
              }
            },
            {
              "output": {
                "data": "123",
                "exist": true
              }
            },
            {
              "output": {
                "data": "123",
                "exist": true
              }
            },
            {
              "output": {
                "data": "123",
                "exist": true
              }
            },
            {
              "output": {
                "data": "123",
                "exist": true
              }
            }
          ],
          "_A.CopyOfDowhile.output": {
            "data": "123",
            "exist": true
          }
        },
        "tasks": {
          "CopyOfDowhile": {
            "change": 1,
            "status": 20,
            "input": {
              "key": "123"
            }
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "24": {
        "change": 2,
        "subflowId": 0,
        "taskId": "CopyOfDowhile"
      },
      "25": {
        "change": 0,
        "subflowId": 0,
        "taskId": "CopyOfDowhile"
      }
    }
  },
  {
    "id": 25,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "CopyOfDowhile",
        "status": 0,
        "attrs": {
          "_A.CopyOfDowhile": [
            {
              "output": {
                "data": "123",
                "exist": true
              }
            },
            {
              "output": {
                "data": "123",
                "exist": true
              }
            },
            {
              "output": {
                "data": "123",
                "exist": true
              }
            },
            {
              "output": {
                "data": "123",
                "exist": true
              }
            },
            {
              "output": {
                "data": "123",
                "exist": true
              }
            },
            {
              "output": {
                "data": "123",
                "exist": true
              }
            }
          ],
          "_A.CopyOfDowhile.output": {
            "data": "123",
            "exist": true
          }
        },
        "tasks": {
          "CopyOfDowhile": {
            "change": 2,
            "status": 40,
            "input": {
              "key": "123"
            }
          },
          "Log": {
            "change": 1,
            "status": 20,
            "input": null
          }
        },
        "links": {
          "3": {
            "change": 2,
            "status": 0,
            "from": "",
            "to": ""
          },
          "4": {
            "change": 1,
            "status": 2,
            "from": "CopyOfDowhile",
            "to": "Log"
          }
        },
        "returnData": null
      }
    },
    "queueChanges": {
      "25": {
        "change": 2,
        "subflowId": 0,
        "taskId": "CopyOfDowhile"
      },
      "26": {
        "change": 0,
        "subflowId": 0,
        "taskId": "Log"
      }
    }
  },
  {
    "id": 26,
    "flowId": "fcd9e55c1d0bb8c3e1f892fbf5f8e032",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "res://flow:flow3",
        "subflowId": 0,
        "taskId": "Log",
        "status": 500,
        "attrs": null,
        "tasks": {
          "Log": {
            "change": 2,
            "status": 40,
            "input": {
              "addDetails": false,
              "message": "aaaaa",
              "usePrint": false
            }
          }
        },
        "links": {
          "4": {
            "change": 2,
            "status": 0,
            "from": "",
            "to": ""
          }
        },
        "returnData": {}
      }
    },
    "queueChanges": {
      "26": {
        "change": 2,
        "subflowId": 0,
        "taskId": "Log"
      }
    }
  }
]`
func TestStepToTask(t *testing.T) {
	var steps []*state.Step

	err := json.Unmarshal([]byte(stepData1), &steps)
	if err != nil {
		panic(err)
	}

	var tasks []*Task
	for _, step := range steps {
		t, err := StepToTaskWithReadyTaskInput(step)
		if err != nil {
			panic(err)
		}
		tasks = append(tasks, t...)
	}

	v, _ := json.Marshal(tasks)
	fmt.Println(string(v))
}

func TestStepToTaskError(t *testing.T) {
	var steps []*state.Step

	err := json.Unmarshal([]byte(stepsData), &steps)
	if err != nil {
		panic(err)
	}

	var tasks []*Task
	for _, step := range steps {
		t, err := StepToTask(step)
		if err != nil {
			panic(err)
		}
		tasks = append(tasks, t...)
	}

	v, _ := json.Marshal(tasks)
	fmt.Println(string(v))
}

var stepsData = `[
  {
    "id": 0,
    "flowId": "c7cfc29deecd3fd1259cba16745b757c",
    "flowChanges": {
      "0": {
        "newFlow": true,
        "flowURI": "res://flow:asdasdasd",
        "subflowId": 0,
        "taskId": "",
        "status": 100,
        "attrs": null,
        "tasks": {
          "LogMessage": {
            "change": 1,
            "status": 20,
            "input": null
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "1": {
        "change": 0,
        "subflowId": 0,
        "taskId": "LogMessage"
      }
    }
  },
  {
    "id": 1,
    "flowId": "c7cfc29deecd3fd1259cba16745b757c",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "",
        "status": 0,
        "attrs": null,
        "tasks": {
          "Log2": {
            "change": 1,
            "status": 20,
            "input": null
          },
          "LogMessage": {
            "change": 2,
            "status": 40,
            "input": {
              "Log Level": "INFO",
              "flowInfo": false,
              "message": "11===1111"
            }
          }
        },
        "links": {
          "0": {
            "change": 1,
            "status": 2,
            "from": "LogMessage",
            "to": "Log2"
          }
        },
        "returnData": null
      }
    },
    "queueChanges": {
      "1": {
        "change": 2,
        "subflowId": 0,
        "taskId": "LogMessage"
      },
      "2": {
        "change": 0,
        "subflowId": 0,
        "taskId": "Log2"
      }
    }
  },
  {
    "id": 2,
    "flowId": "c7cfc29deecd3fd1259cba16745b757c",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "",
        "status": 0,
        "attrs": null,
        "tasks": {
          "Log2": {
            "change": 2,
            "status": 40,
            "input": {
              "Log Level": "INFO",
              "flowInfo": false,
              "message": "22===22"
            }
          },
          "Log3": {
            "change": 1,
            "status": 20,
            "input": null
          }
        },
        "links": {
          "0": {
            "change": 2,
            "status": 2,
            "from": "Log2",
            "to": "Log3"
          }
        },
        "returnData": null
      }
    },
    "queueChanges": {
      "2": {
        "change": 2,
        "subflowId": 0,
        "taskId": "Log2"
      },
      "3": {
        "change": 0,
        "subflowId": 0,
        "taskId": "Log3"
      }
    }
  },
  {
    "id": 3,
    "flowId": "c7cfc29deecd3fd1259cba16745b757c",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "",
        "status": 0,
        "attrs": null,
        "tasks": {
          "Log3": {
            "change": 2,
            "status": 40,
            "input": {
              "Log Level": "INFO",
              "flowInfo": false,
              "message": "33=33"
            }
          },
          "Log4": {
            "change": 1,
            "status": 20,
            "input": null
          }
        },
        "links": {
          "0": {
            "change": 1,
            "status": 2,
            "from": "Log3",
            "to": "Log4"
          },
          "1": {
            "change": 2,
            "status": 0,
            "from": "",
            "to": ""
          }
        },
        "returnData": null
      }
    },
    "queueChanges": {
      "3": {
        "change": 2,
        "subflowId": 0,
        "taskId": "Log3"
      },
      "4": {
        "change": 0,
        "subflowId": 0,
        "taskId": "Log4"
      }
    }
  },
  {
    "id": 4,
    "flowId": "c7cfc29deecd3fd1259cba16745b757c",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "",
        "status": 0,
        "attrs": null,
        "tasks": {
          "Log4": {
            "change": 2,
            "status": 40,
            "input": {
              "Log Level": "INFO",
              "flowInfo": false,
              "message": "44=444"
            }
          },
          "Log5": {
            "change": 1,
            "status": 20,
            "input": null
          }
        },
        "links": {
          "0": {
            "change": 1,
            "status": 2,
            "from": "Log4",
            "to": "Log5"
          },
          "2": {
            "change": 2,
            "status": 0,
            "from": "",
            "to": ""
          }
        },
        "returnData": null
      }
    },
    "queueChanges": {
      "4": {
        "change": 2,
        "subflowId": 0,
        "taskId": "Log4"
      },
      "5": {
        "change": 0,
        "subflowId": 0,
        "taskId": "Log5"
      }
    }
  },
  {
    "id": 5,
    "flowId": "c7cfc29deecd3fd1259cba16745b757c",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "",
        "status": 0,
        "attrs": null,
        "tasks": {
          "Log5": {
            "change": 2,
            "status": 40,
            "input": {
              "Log Level": "INFO",
              "flowInfo": false,
              "message": "5==55555"
            }
          },
          "Log6": {
            "change": 1,
            "status": 20,
            "input": null
          }
        },
        "links": {
          "0": {
            "change": 1,
            "status": 2,
            "from": "Log5",
            "to": "Log6"
          },
          "3": {
            "change": 2,
            "status": 0,
            "from": "",
            "to": ""
          }
        },
        "returnData": null
      }
    },
    "queueChanges": {
      "5": {
        "change": 2,
        "subflowId": 0,
        "taskId": "Log5"
      },
      "6": {
        "change": 0,
        "subflowId": 0,
        "taskId": "Log6"
      }
    }
  },
  {
    "id": 6,
    "flowId": "c7cfc29deecd3fd1259cba16745b757c",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "",
        "status": 0,
        "attrs": null,
        "tasks": {
          "Log6": {
            "change": 2,
            "status": 40,
            "input": {
              "Log Level": "INFO",
              "flowInfo": false,
              "message": "66===666"
            }
          },
          "StartaSubFlow": {
            "change": 1,
            "status": 20,
            "input": null
          }
        },
        "links": {
          "0": {
            "change": 1,
            "status": 2,
            "from": "Log6",
            "to": "StartaSubFlow"
          },
          "4": {
            "change": 2,
            "status": 0,
            "from": "",
            "to": ""
          }
        },
        "returnData": null
      }
    },
    "queueChanges": {
      "6": {
        "change": 2,
        "subflowId": 0,
        "taskId": "Log6"
      },
      "7": {
        "change": 0,
        "subflowId": 0,
        "taskId": "StartaSubFlow"
      }
    }
  },
  {
    "id": 7,
    "flowId": "c7cfc29deecd3fd1259cba16745b757c",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "StartaSubFlow",
        "status": 0,
        "attrs": null,
        "tasks": {
          "StartaSubFlow": {
            "change": 1,
            "status": 30,
            "input": {
              "input": "INPUT"
            }
          }
        },
        "links": null,
        "returnData": null
      },
      "1": {
        "newFlow": true,
        "flowURI": "res://flow:subflow",
        "subflowId": 1,
        "taskId": "StartaSubFlow",
        "status": 100,
        "attrs": {
          "input": "INPUT"
        },
        "tasks": {
          "LogMessage": {
            "change": 1,
            "status": 20,
            "input": null
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "7": {
        "change": 2,
        "subflowId": 0,
        "taskId": "StartaSubFlow"
      },
      "8": {
        "change": 0,
        "subflowId": 0,
        "taskId": "LogMessage"
      }
    }
  },
  {
    "id": 8,
    "flowId": "c7cfc29deecd3fd1259cba16745b757c",
    "flowChanges": {
      "1": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "",
        "status": 0,
        "attrs": null,
        "tasks": {
          "InvokeRESTService": {
            "change": 1,
            "status": 20,
            "input": null
          },
          "LogMessage": {
            "change": 2,
            "status": 40,
            "input": {
              "Log Level": "INFO",
              "flowInfo": false,
              "message": "subflow log message"
            }
          }
        },
        "links": {
          "0": {
            "change": 1,
            "status": 2,
            "from": "LogMessage",
            "to": "InvokeRESTService"
          }
        },
        "returnData": null
      }
    },
    "queueChanges": {
      "8": {
        "change": 2,
        "subflowId": 0,
        "taskId": "LogMessage"
      },
      "9": {
        "change": 0,
        "subflowId": 0,
        "taskId": "InvokeRESTService"
      }
    }
  },
  {
    "id": 9,
    "flowId": "c7cfc29deecd3fd1259cba16745b757c",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "StartaSubFlow",
        "status": 0,
        "attrs": {
          "_E": {
            "activity": "InvokeRESTService",
            "code": "",
            "data": null,
            "message": "Get \"\": unsupported protocol scheme \"\"",
            "type": "activity"
          },
          "_E.StartaSubFlow": {
            "activity": "InvokeRESTService",
            "code": "",
            "data": null,
            "message": "Get \"\": unsupported protocol scheme \"\"",
            "type": "activity"
          }
        },
        "tasks": {
          "LogMessage1": {
            "change": 1,
            "status": 20,
            "input": null
          },
          "StartaSubFlow": {
            "change": 1,
            "status": 100,
            "input": {
              "input": "INPUT"
            }
          }
        },
        "links": null,
        "returnData": null
      },
      "1": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "InvokeRESTService",
        "status": 700,
        "attrs": {
          "_E": {
            "activity": "InvokeRESTService",
            "code": "",
            "data": null,
            "message": "Get \"\": unsupported protocol scheme \"\"",
            "type": "activity"
          },
          "_E.InvokeRESTService": {
            "activity": "InvokeRESTService",
            "code": "",
            "data": null,
            "message": "Get \"\": unsupported protocol scheme \"\"",
            "type": "activity"
          }
        },
        "tasks": {
          "InvokeRESTService": {
            "change": 1,
            "status": 100,
            "input": {
              "Method": "GET",
              "Timeout": 0,
              "Use certificate for verification": false,
              "enableASR": false,
              "requestType": "application/json"
            }
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "10": {
        "change": 0,
        "subflowId": 0,
        "taskId": "LogMessage1"
      },
      "9": {
        "change": 2,
        "subflowId": 0,
        "taskId": "InvokeRESTService"
      }
    }
  },
  {
    "id": 10,
    "flowId": "c7cfc29deecd3fd1259cba16745b757c",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "",
        "status": 0,
        "attrs": null,
        "tasks": {
          "InvokeRESTService": {
            "change": 1,
            "status": 20,
            "input": null
          },
          "LogMessage1": {
            "change": 2,
            "status": 40,
            "input": {
              "Log Level": "INFO",
              "flowInfo": false,
              "message": "ERROR HANDLER MESSAGE:Get \"\": unsupported protocol scheme \"\""
            }
          }
        },
        "links": {
          "0": {
            "change": 1,
            "status": 2,
            "from": "LogMessage1",
            "to": "InvokeRESTService"
          }
        },
        "returnData": null
      }
    },
    "queueChanges": {
      "10": {
        "change": 2,
        "subflowId": 0,
        "taskId": "LogMessage1"
      },
      "11": {
        "change": 0,
        "subflowId": 0,
        "taskId": "InvokeRESTService"
      }
    }
  },
  {
    "id": 11,
    "flowId": "c7cfc29deecd3fd1259cba16745b757c",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "InvokeRESTService",
        "status": 700,
        "attrs": null,
        "tasks": {
          "InvokeRESTService": {
            "change": 1,
            "status": 100,
            "input": {
              "Method": "GET",
              "Timeout": 0,
              "Use certificate for verification": false,
              "enableASR": false,
              "requestType": "application/json"
            }
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "11": {
        "change": 2,
        "subflowId": 0,
        "taskId": "InvokeRESTService"
      }
    }
  }
]`

var stepsData2 = `[{"id":0,"flowId":"58cfc1b3b2a49d1960b6ed47bc2834cd","flowChanges":{"0":{"newFlow":true,"flowURI":"res://flow:asdasdasd","subflowId":0,"taskId":"","status":100,"attrs":null,"tasks":{"LogMessage":{"change":1,"status":20,"input":null}},"links":null,"returnData":null}},"queueChanges":{"1":{"change":0,"subflowId":0,"taskId":"LogMessage"}}},{"id":1,"flowId":"58cfc1b3b2a49d1960b6ed47bc2834cd","flowChanges":{"0":{"newFlow":false,"flowURI":"","subflowId":0,"taskId":"","status":0,"attrs":null,"tasks":{"Log2":{"change":1,"status":20,"input":null},"LogMessage":{"change":2,"status":40,"input":{"Log Level":"INFO","flowInfo":false,"message":"11===1111"}}},"links":{"0":{"change":1,"status":2,"from":"LogMessage","to":"Log2"}},"returnData":null}},"queueChanges":{"1":{"change":2,"subflowId":0,"taskId":"LogMessage"},"2":{"change":0,"subflowId":0,"taskId":"Log2"}}},{"id":2,"flowId":"58cfc1b3b2a49d1960b6ed47bc2834cd","flowChanges":{"0":{"newFlow":false,"flowURI":"","subflowId":0,"taskId":"","status":0,"attrs":null,"tasks":{"Log2":{"change":2,"status":40,"input":{"Log Level":"INFO","flowInfo":false,"message":"22===22"}},"Log3":{"change":1,"status":20,"input":null}},"links":{"0":{"change":2,"status":2,"from":"Log2","to":"Log3"}},"returnData":null}},"queueChanges":{"2":{"change":2,"subflowId":0,"taskId":"Log2"},"3":{"change":0,"subflowId":0,"taskId":"Log3"}}},{"id":3,"flowId":"58cfc1b3b2a49d1960b6ed47bc2834cd","flowChanges":{"0":{"newFlow":false,"flowURI":"","subflowId":0,"taskId":"","status":0,"attrs":null,"tasks":{"Log3":{"change":2,"status":40,"input":{"Log Level":"INFO","flowInfo":false,"message":"33=33"}},"Log4":{"change":1,"status":20,"input":null}},"links":{"0":{"change":1,"status":2,"from":"Log3","to":"Log4"},"1":{"change":2,"status":0,"from":"","to":""}},"returnData":null}},"queueChanges":{"3":{"change":2,"subflowId":0,"taskId":"Log3"},"4":{"change":0,"subflowId":0,"taskId":"Log4"}}},{"id":4,"flowId":"58cfc1b3b2a49d1960b6ed47bc2834cd","flowChanges":{"0":{"newFlow":false,"flowURI":"","subflowId":0,"taskId":"","status":0,"attrs":null,"tasks":{"Log4":{"change":2,"status":40,"input":{"Log Level":"INFO","flowInfo":false,"message":"44=444"}},"Log5":{"change":1,"status":20,"input":null}},"links":{"0":{"change":1,"status":2,"from":"Log4","to":"Log5"},"2":{"change":2,"status":0,"from":"","to":""}},"returnData":null}},"queueChanges":{"4":{"change":2,"subflowId":0,"taskId":"Log4"},"5":{"change":0,"subflowId":0,"taskId":"Log5"}}},{"id":5,"flowId":"58cfc1b3b2a49d1960b6ed47bc2834cd","flowChanges":{"0":{"newFlow":false,"flowURI":"","subflowId":0,"taskId":"","status":0,"attrs":null,"tasks":{"Log5":{"change":2,"status":40,"input":{"Log Level":"INFO","flowInfo":false,"message":"5==55555"}},"Log6":{"change":1,"status":20,"input":null}},"links":{"0":{"change":1,"status":2,"from":"Log5","to":"Log6"},"3":{"change":2,"status":0,"from":"","to":""}},"returnData":null}},"queueChanges":{"5":{"change":2,"subflowId":0,"taskId":"Log5"},"6":{"change":0,"subflowId":0,"taskId":"Log6"}}},{"id":6,"flowId":"58cfc1b3b2a49d1960b6ed47bc2834cd","flowChanges":{"0":{"newFlow":false,"flowURI":"","subflowId":0,"taskId":"","status":0,"attrs":null,"tasks":{"Log6":{"change":2,"status":40,"input":{"Log Level":"INFO","flowInfo":false,"message":"66===666"}},"StartaSubFlow":{"change":1,"status":20,"input":null}},"links":{"0":{"change":1,"status":2,"from":"Log6","to":"StartaSubFlow"},"4":{"change":2,"status":0,"from":"","to":""}},"returnData":null}},"queueChanges":{"6":{"change":2,"subflowId":0,"taskId":"Log6"},"7":{"change":0,"subflowId":0,"taskId":"StartaSubFlow"}}},{"id":7,"flowId":"58cfc1b3b2a49d1960b6ed47bc2834cd","flowChanges":{"0":{"newFlow":false,"flowURI":"","subflowId":0,"taskId":"StartaSubFlow","status":0,"attrs":null,"tasks":{"StartaSubFlow":{"change":1,"status":30,"input":{"input":"INPUT"}}},"links":null,"returnData":null},"1":{"newFlow":true,"flowURI":"res://flow:subflow","subflowId":1,"taskId":"StartaSubFlow","status":100,"attrs":{"input":"INPUT"},"tasks":{"LogMessage":{"change":1,"status":20,"input":null}},"links":null,"returnData":null}},"queueChanges":{"7":{"change":2,"subflowId":0,"taskId":"StartaSubFlow"},"8":{"change":0,"subflowId":0,"taskId":"LogMessage"}}},{"id":8,"flowId":"58cfc1b3b2a49d1960b6ed47bc2834cd","flowChanges":{"1":{"newFlow":false,"flowURI":"","subflowId":0,"taskId":"","status":0,"attrs":null,"tasks":{"InvokeRESTService":{"change":1,"status":20,"input":null},"LogMessage":{"change":2,"status":40,"input":{"Log Level":"INFO","flowInfo":false,"message":"subflow log message"}}},"links":{"0":{"change":1,"status":2,"from":"LogMessage","to":"InvokeRESTService"}},"returnData":null}},"queueChanges":{"8":{"change":2,"subflowId":0,"taskId":"LogMessage"},"9":{"change":0,"subflowId":0,"taskId":"InvokeRESTService"}}},{"id":9,"flowId":"58cfc1b3b2a49d1960b6ed47bc2834cd","flowChanges":{"1":{"newFlow":false,"flowURI":"","subflowId":0,"taskId":"InvokeRESTService","status":0,"attrs":{"_E":{"activity":"InvokeRESTService","code":"","data":null,"message":"Get \"\": unsupported protocol scheme \"\"","type":"activity"},"_E.InvokeRESTService":{"activity":"InvokeRESTService","code":"","data":null,"message":"Get \"\": unsupported protocol scheme \"\"","type":"activity"}},"tasks":{"InvokeRESTService":{"change":1,"status":100,"input":{"Method":"GET","Timeout":0,"Use certificate for verification":false,"enableASR":false,"requestType":"application/json"}},"LogMessage1":{"change":1,"status":20,"input":null}},"links":null,"returnData":null}},"queueChanges":{"10":{"change":0,"subflowId":0,"taskId":"LogMessage1"},"9":{"change":2,"subflowId":0,"taskId":"InvokeRESTService"}}},{"id":10,"flowId":"58cfc1b3b2a49d1960b6ed47bc2834cd","flowChanges":{"1":{"newFlow":false,"flowURI":"","subflowId":0,"taskId":"","status":0,"attrs":null,"tasks":{"LogMessage1":{"change":2,"status":40,"input":{"Log Level":"INFO","flowInfo":false,"message":"Newsubflowerror handler  error: Get \"\": unsupported protocol scheme \"\""}},"Return":{"change":1,"status":20,"input":null}},"links":{"0":{"change":1,"status":2,"from":"LogMessage1","to":"Return"}},"returnData":null}},"queueChanges":{"10":{"change":2,"subflowId":0,"taskId":"LogMessage1"},"11":{"change":0,"subflowId":0,"taskId":"Return"}}},{"id":11,"flowId":"58cfc1b3b2a49d1960b6ed47bc2834cd","flowChanges":{"1":{"newFlow":false,"flowURI":"res://flow:subflow","subflowId":0,"taskId":"Return","status":500,"attrs":null,"tasks":{"Return":{"change":2,"status":40,"input":null}},"links":{"1":{"change":2,"status":0,"from":"","to":""}},"returnData":{"output":"dasdasdasdasda"}}},"queueChanges":{"11":{"change":2,"subflowId":0,"taskId":"Return"},"12":{"change":0,"subflowId":0,"taskId":"StartaSubFlow"}}},{"id":12,"flowId":"58cfc1b3b2a49d1960b6ed47bc2834cd","flowChanges":{"0":{"newFlow":false,"flowURI":"","subflowId":0,"taskId":"","status":0,"attrs":{"_A.StartaSubFlow.output":"dasdasdasdasda"},"tasks":{"StartaSubFlow":{"change":2,"status":40,"input":{"input":"INPUT"}},"ThrowError":{"change":1,"status":20,"input":null}},"links":{"0":{"change":1,"status":2,"from":"StartaSubFlow","to":"ThrowError"},"5":{"change":2,"status":0,"from":"","to":""}},"returnData":null}},"queueChanges":{"12":{"change":2,"subflowId":0,"taskId":"StartaSubFlow"},"13":{"change":0,"subflowId":0,"taskId":"ThrowError"}}},{"id":13,"flowId":"58cfc1b3b2a49d1960b6ed47bc2834cd","flowChanges":{"0":{"newFlow":false,"flowURI":"","subflowId":0,"taskId":"ThrowError","status":0,"attrs":{"_E":{"activity":"ThrowError","code":"","data":null,"message":"main flow throw error","type":"activity"},"_E.ThrowError":{"activity":"ThrowError","code":"","data":null,"message":"main flow throw error","type":"activity"}},"tasks":{"LogMessage1":{"change":1,"status":20,"input":null},"ThrowError":{"change":1,"status":100,"input":{"message":"main flow throw error"}}},"links":null,"returnData":null}},"queueChanges":{"13":{"change":2,"subflowId":0,"taskId":"ThrowError"},"14":{"change":0,"subflowId":0,"taskId":"LogMessage1"}}},{"id":14,"flowId":"58cfc1b3b2a49d1960b6ed47bc2834cd","flowChanges":{"0":{"newFlow":false,"flowURI":"","subflowId":0,"taskId":"","status":0,"attrs":null,"tasks":{"InvokeRESTService":{"change":1,"status":20,"input":null},"LogMessage1":{"change":2,"status":40,"input":{"Log Level":"INFO","flowInfo":false,"message":"ERROR HANDLER MESSAGE:main flow throw error"}}},"links":{"0":{"change":1,"status":2,"from":"LogMessage1","to":"InvokeRESTService"}},"returnData":null}},"queueChanges":{"14":{"change":2,"subflowId":0,"taskId":"LogMessage1"},"15":{"change":0,"subflowId":0,"taskId":"InvokeRESTService"}}},{"id":15,"flowId":"58cfc1b3b2a49d1960b6ed47bc2834cd","flowChanges":{"0":{"newFlow":false,"flowURI":"","subflowId":0,"taskId":"InvokeRESTService","status":700,"attrs":null,"tasks":{"InvokeRESTService":{"change":1,"status":100,"input":{"Method":"GET","Timeout":0,"Use certificate for verification":false,"enableASR":false,"requestType":"application/json"}}},"links":null,"returnData":null}},"queueChanges":{"15":{"change":2,"subflowId":0,"taskId":"InvokeRESTService"}}}]
`
var subflowStartStep = `{
    "id": 2,
    "flowId": "0c8b85f5e41b0766c5f765ed573d23b0",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "StartaSubFlow",
        "status": 0,
        "attrs": null,
        "tasks": {
          "StartaSubFlow": {
            "change": 1,
            "status": 30,
            "input": {
              "input": ""
            }
          }
        },
        "links": null,
        "returnData": null
      },
      "1": {
        "newFlow": true,
        "flowURI": "res://flow:subflow",
        "subflowId": 1,
        "taskId": "LogMessage",
        "status": 100,
        "attrs": {
          "input": ""
        },
        "tasks": {
          "LogMessage": {
            "change": 1,
            "status": 20,
            "input": null
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "2": {
        "change": 2,
        "subflowId": 0,
        "taskId": "StartaSubFlow"
      },
      "3": {
        "change": 0,
        "subflowId": 0,
        "taskId": "LogMessage"
      }
    }
  }`

var subflowErrorSteps = `{
    "id": 4,
    "flowId": "0c8b85f5e41b0766c5f765ed573d23b0",
    "flowChanges": {
      "0": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "StartaSubFlow",
        "status": 0,
        "attrs": {
          "_E.StartaSubFlow": {
            "activity": "ThrowError",
            "code": "",
            "data": null,
            "message": "subflow throw error",
            "type": "activity"
          }
        },
        "tasks": {
          "InvokeRESTService": {
            "change": 1,
            "status": 50,
            "input": null
          },
          "LogMessage1": {
            "change": 2,
            "status": 50,
            "input": null
          },
          "LogMessage2": {
            "change": 1,
            "status": 20,
            "input": null
          },
          "Return": {
            "change": 1,
            "status": 50,
            "input": null
          },
          "StartaSubFlow": {
            "change": 2,
            "status": 100,
            "input": {
              "input": ""
            }
          }
        },
        "links": {
          "1": {
            "change": 2,
            "status": 1,
            "from": "StartaSubFlow",
            "to": "LogMessage1"
          },
          "2": {
            "change": 1,
            "status": 3,
            "from": "LogMessage1",
            "to": "InvokeRESTService"
          },
          "3": {
            "change": 1,
            "status": 3,
            "from": "InvokeRESTService",
            "to": "Return"
          },
          "4": {
            "change": 1,
            "status": 2,
            "from": "StartaSubFlow",
            "to": "LogMessage2"
          }
        },
        "returnData": null
      },
      "1": {
        "newFlow": false,
        "flowURI": "",
        "subflowId": 0,
        "taskId": "ThrowError",
        "status": 700,
        "attrs": {
          "_E": {
            "activity": "ThrowError",
            "code": "",
            "data": null,
            "message": "subflow throw error",
            "type": "activity"
          },
          "_E.ThrowError": {
            "activity": "ThrowError",
            "code": "",
            "data": null,
            "message": "subflow throw error",
            "type": "activity"
          }
        },
        "tasks": {
          "ThrowError": {
            "change": 1,
            "status": 100,
            "input": {
              "message": "subflow throw error"
            }
          }
        },
        "links": null,
        "returnData": null
      }
    },
    "queueChanges": {
      "4": {
        "change": 2,
        "subflowId": 0,
        "taskId": "ThrowError"
      },
      "5": {
        "change": 0,
        "subflowId": 0,
        "taskId": "LogMessage2"
      }
    }
  }`

func TestTaskSubflowStart(t *testing.T) {
	var step *state.Step
	err := json.Unmarshal([]byte(subflowStartStep), &step)
	if err != nil {
		t.Fatal(err)
	}
	task, err := StepToTask(step)
	if err != nil {
		t.Fatal(err)
	}
	v, _ := json.Marshal(task)
	fmt.Println(string(v))
}

func TestTaskSubflowError(t *testing.T) {
	var step *state.Step
	err := json.Unmarshal([]byte(subflowErrorSteps), &step)
	if err != nil {
		t.Fatal(err)
	}
	task, err := StepToTask(step)
	if err != nil {
		t.Fatal(err)
	}
	v, _ := json.Marshal(task)
	fmt.Println(string(v))
}

var testData = `[
      {
        "id": 0,
        "flowId": "2e2842d9a30937083491f055d89fa1ec",
        "flowChanges": {
          "0": {
            "newFlow": true,
            "flowURI": "res://flow:Test2",
            "subflowId": 0,
            "taskId": "LogMessage",
            "status": 100,
            "attrs": {
              "username": "Tracy"
            },
            "tasks": {
              "LogMessage": {
                "change": 1,
                "status": 20,
                "input": null
              }
            },
            "links": null,
            "returnData": null
          }
        },
        "queueChanges": {
          "1": {
            "change": 0,
            "subflowId": 0,
            "taskId": "LogMessage"
          }
        }
      },
      {
        "id": 1,
        "flowId": "2e2842d9a30937083491f055d89fa1ec",
        "flowChanges": {
          "0": {
            "newFlow": false,
            "flowURI": "",
            "subflowId": 0,
            "taskId": "LogMessage",
            "status": 0,
            "attrs": null,
            "tasks": {
              "InvokeRESTService": {
                "change": 1,
                "status": 20,
                "input": null
              },
              "LogMessage": {
                "change": 2,
                "status": 40,
                "input": {
                  "Log Level": "INFO",
                  "flowInfo": false,
                  "message": "Hello Tracy!! Starting main flow"
                }
              }
            },
            "links": {
              "0": {
                "change": 1,
                "status": 2,
                "from": "LogMessage",
                "to": "InvokeRESTService"
              }
            },
            "returnData": null
          }
        },
        "queueChanges": {
          "1": {
            "change": 2,
            "subflowId": 0,
            "taskId": "LogMessage"
          },
          "2": {
            "change": 0,
            "subflowId": 0,
            "taskId": "InvokeRESTService"
          }
        }
      },
      {
        "id": 2,
        "flowId": "2e2842d9a30937083491f055d89fa1ec",
        "flowChanges": {
          "0": {
            "newFlow": false,
            "flowURI": "",
            "subflowId": 0,
            "taskId": "InvokeRESTService",
            "status": 0,
            "attrs": {
              "_A.InvokeRESTService.configureResponseCodes": false,
              "_A.InvokeRESTService.headers": {
                "Content-Type": "application/json"
              },
              "_A.InvokeRESTService.responseBody": [
                {
                  "category": {
                    "id": 523,
                    "name": "gctIeDFo5_1LbSJ8"
                  },
                  "id": 69,
                  "name": "doggie",
                  "photoUrls": [
                    "IZ2cDIXdIuexxc9T"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 777,
                      "name": "8I6x_I8xG_YyO53w; ping -n 4 localhost"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 434,
                    "name": "FL8STnOnFE8pYgl2"
                  },
                  "id": 557,
                  "name": "doggie",
                  "photoUrls": [
                    "eLz8pw8yke2j-l_2"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 998,
                      "name": "KmJ3KWW34vxcT7W5 AND 1 IN (SELECT BENCHMARK(2*15000000,MD5(CHAR(97))))-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 4,
                    "name": "1HU9zwcdiX-17u7f"
                  },
                  "id": 904,
                  "name": "doggie",
                  "photoUrls": [
                    "w6PYmahm0YX3NRBu"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 449,
                      "name": "xSSey2yoF7az3M04 AND 1 IN (SELECT BENCHMARK(4*15000000,MD5(CHAR(97))))-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 172,
                    "name": "CRoJ8dmozs8daLpE"
                  },
                  "id": 156,
                  "name": "doggie",
                  "photoUrls": [
                    "46OJpz-Buul-O0pR"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 667,
                      "name": "tUhcy9BgGJukQnqY AND 1 IN (SELECT BENCHMARK(8*15000000,MD5(CHAR(97))))-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 741,
                    "name": "rOB_UiutVtktgXsH"
                  },
                  "id": 621,
                  "name": "doggie",
                  "photoUrls": [
                    "ZfpGlZkuB6hF9cV6"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 664,
                      "name": "4WuCb3jud3l9AHMU' AND 1 IN (SELECT BENCHMARK(8*15000000,MD5(CHAR(97))))-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 981,
                    "name": "B0t0tcF3f6SQICx6"
                  },
                  "id": 911,
                  "name": "doggie",
                  "photoUrls": [
                    "TmSkbUCoyCdNXqLy"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 903,
                      "name": "9dWyx9cOMIDqV8JO)) OR SLEEP(2)=0 LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 310,
                    "name": "1MP0--oe33OkvxfD"
                  },
                  "id": 491,
                  "name": "doggie",
                  "photoUrls": [
                    "PU9mF_ptB7PSFzha"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 483,
                      "name": "s2Mzx6mXou7s0w6G=SLEEP(4)="
                    }
                  ]
                },
                {
                  "category": {
                    "id": 564,
                    "name": "KrSiGFlUhIs_pNEM"
                  },
                  "id": 60,
                  "name": "doggie",
                  "photoUrls": [
                    "8_HYJ6acUHYmhSpW"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 578,
                      "name": "DNiJtIKgS6UB8BSF&& sleep 2"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 841,
                    "name": "CVqXoe_AgeHWG2P_"
                  },
                  "id": 444,
                  "name": "doggie",
                  "photoUrls": [
                    "W3XgFBxkCykorrLA"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 377,
                      "name": "8KMUs93z_soAz0Ro' WHERE SLEEP(2) LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 451,
                    "name": "bjGsZDs135w9IpWp"
                  },
                  "id": 519,
                  "name": "doggie",
                  "photoUrls": [
                    "KvSAR4oSV7ECWGua"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 197,
                      "name": "1Zs3csWbZQOlrA5f'));WAITFOR DELAY '00:00:2'-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 762,
                    "name": "TSKmgPvB0Gv1dpZo"
                  },
                  "id": 589,
                  "name": "doggie",
                  "photoUrls": [
                    "ilU6DWuE3tJo0f48"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 513,
                      "name": "_O5TUszzRxkWs9e2' WHERE SLEEP(2) LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 740,
                    "name": "whXHstDP5DsDdxt1"
                  },
                  "id": 448,
                  "name": "doggie",
                  "photoUrls": [
                    "THRAGKr8KfZbDQUi"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 711,
                      "name": "481ExghT-f4J9eAB OR SLEEP(8)=0 LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 902,
                    "name": "LRDEcwYSl-NQj0Sy"
                  },
                  "id": 713,
                  "name": "doggie",
                  "photoUrls": [
                    "c0UKWbS7F0qH4XZV"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 540,
                      "name": "5SserYx3TK8K3TRH AND SLEEP(8)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 485,
                    "name": "aMz08LZIvmo1Mryd"
                  },
                  "id": 469,
                  "name": "doggie",
                  "photoUrls": [
                    "QQaEygO0NFKQ1HVZ"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 197,
                      "name": "--kqUxxYq5CBmzhF') AND SLEEP(2)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 516,
                    "name": "D6mD1zgdCNEWgOqS"
                  },
                  "id": 814,
                  "name": "doggie",
                  "photoUrls": [
                    "-xAokqOP5OKB88ye"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 715,
                      "name": "A1CncQzl4--zwCeK' AND SLEEP(4)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 327,
                    "name": "YCI75LhrNPCgtNXE"
                  },
                  "id": 756,
                  "name": "doggie",
                  "photoUrls": [
                    "3omjjUBNruCNgCiZ"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 555,
                      "name": "B2uuCHYGXaTSDHlw;select pg_sleep(8); -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 357,
                    "name": "nwBFoXkD_W6qtupF"
                  },
                  "id": 950,
                  "name": "doggie",
                  "photoUrls": [
                    "64owgvmBsP8aByQ1"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 210,
                      "name": "7pW8PBceaHtoWq15' OR SLEEP(4)=0 LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 566,
                    "name": "y6sSCksQzOlj_FQP"
                  },
                  "id": 669,
                  "name": "doggie",
                  "photoUrls": [
                    "T1x1oY1zrJ-3e0S9"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 268,
                      "name": "M6K5aecr4UEoHKc5';select pg_sleep(8); -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 920,
                    "name": "LiqsTaKE-OJQbnPJ"
                  },
                  "id": 975,
                  "name": "doggie",
                  "photoUrls": [
                    "93qiRGgj_OUd-53j"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 508,
                      "name": "Z1pQs4jaKo46IOnx' AND SLEEP(4)=0 LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 161,
                    "name": "3rRDZI9nogQxDnnv"
                  },
                  "id": 30,
                  "name": "doggie",
                  "photoUrls": [
                    "TH0uy4fdgzo4aJ4l"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 753,
                      "name": "UR7ziabZACenOyBn' AND SLEEP(2)=0 LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 843,
                    "name": "RYYnzCgjwgdGGBd-"
                  },
                  "id": 481,
                  "name": "doggie",
                  "photoUrls": [
                    "W9H5Zh3piEDuse0S"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 1024,
                      "name": "PgT2NXTHjxYvOGq1"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 396,
                    "name": "Zm_e5DGFgM8SqkTR"
                  },
                  "id": 853,
                  "name": "doggie",
                  "photoUrls": [
                    "z0AaJb-J1KT4oTpC"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 108,
                      "name": "A5f0Om0jh4DdMfPy' AND SLEEP(2)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 687,
                    "name": "W900H0yol_4t0mqC"
                  },
                  "id": 307,
                  "name": "doggie",
                  "photoUrls": [
                    "jaCxWqXc-ABMM_wo"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 942,
                      "name": "W6zxBndZ9oXeBF91)) AND SLEEP(4)=0 LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 46,
                    "name": "QO3SQ0ewFTR0Uuha"
                  },
                  "id": 207,
                  "name": "doggie",
                  "photoUrls": [
                    "6Uki1wHGvQcJ2Av0"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 91,
                      "name": "dtjTn_ZoDS3d8wW-& ping -n 4 localhost"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 939,
                    "name": "qw7QozTf5NYw4GHh"
                  },
                  "id": 869,
                  "name": "doggie",
                  "photoUrls": [
                    "3HMneDWlG-CZtpe-"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 320,
                      "name": "wp2ZuKed4DJTxb3Y) OR SLEEP(8)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 12,
                    "name": "7GGANOCEq8o2u_o-"
                  },
                  "id": 697,
                  "name": "doggie",
                  "photoUrls": [
                    "BRH7iNYhZHGQ0XMx"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 356,
                      "name": "4_ZtnpkRHihfSaKp' OR SLEEP(2)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 752,
                    "name": "AZ3U-9ErwQJGcqF8"
                  },
                  "id": 148,
                  "name": "doggie",
                  "photoUrls": [
                    "cRJYCR75ispHdSD8"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 992,
                      "name": "Upx9truuU-3nA7K_') OR SLEEP(8)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 984,
                    "name": "X2j7c35qmtmWiu8q"
                  },
                  "id": 163,
                  "name": "doggie",
                  "photoUrls": [
                    "GEiVu5GrlMyMUA0y"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 224,
                      "name": "kJfy2qbP3nLwmQaB') OR SLEEP(4)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 692,
                    "name": "M6eGOGs-UBQkFNNI"
                  },
                  "id": 726,
                  "name": "doggie",
                  "photoUrls": [
                    "ss9I7ZpID2opmBB4"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 705,
                      "name": "0be7DX_gdB1fM79l') OR SLEEP(4)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 412,
                    "name": "7zy8UGAeRD9xjlA6"
                  },
                  "id": 326,
                  "name": "doggie",
                  "photoUrls": [
                    "Cc80cEVFBi9zqEE5"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 447,
                      "name": "_K8Uvrq38r0-Lyw2' AND 1 IN (SELECT BENCHMARK(8*15000000,MD5(CHAR(97))))-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 1009,
                    "name": "1JmzJE_ozQv2jIgo"
                  },
                  "id": 900,
                  "name": "doggie",
                  "photoUrls": [
                    "543sjjikvNYjn0As"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 853,
                      "name": "TX-bdtrlEs3CkrIX AND SLEEP(8)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 876,
                    "name": "3itvGTIydY7euNqh"
                  },
                  "id": 301,
                  "name": "doggie",
                  "photoUrls": [
                    "WJp2Bmx4PY2fpRw3"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 631,
                      "name": "sf7Js1uNEk-0kNxF) AND SLEEP(4)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 636,
                    "name": "5X_j-XzuOaNSzuQp"
                  },
                  "id": 637,
                  "name": "doggie",
                  "photoUrls": [
                    "C9CRzNPzytjRL5Dj"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 302,
                      "name": "ZNzUbCO-U8MzhgHN WHERE SLEEP(2) LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 753,
                    "name": "jScic3RyV5LyzK28"
                  },
                  "id": 284,
                  "name": "doggie",
                  "photoUrls": [
                    "3nwXeEjCbG5kc61C"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 170,
                      "name": "Bdb5y8bLIY6ABmtU' AND SLEEP(2)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 233,
                    "name": "C0bYq_NeNiSBnHSd"
                  },
                  "id": 917,
                  "name": "doggie",
                  "photoUrls": [
                    "s5OCAeq6av1-AvMm"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 459,
                      "name": "GG2Id6usMZjvNiww' AND SLEEP(4)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 764,
                    "name": "To5SyGaJ50a8qU5_"
                  },
                  "id": 721,
                  "name": "doggie",
                  "photoUrls": [
                    "JGqmMwp5Zst30H58"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 282,
                      "name": "LfQ7-cZAHqEQ2Xpy) AND SLEEP(2)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 521,
                    "name": "qwdJsZ71kvRXWybb"
                  },
                  "id": 748,
                  "name": "doggie",
                  "photoUrls": [
                    "A1Ir4EInla82hb_f"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 752,
                      "name": "Tya72szPzzZ_bYet&& sleep 2"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 178,
                    "name": "JMaeCaeLKEAnvyvp"
                  },
                  "id": 740,
                  "name": "doggie",
                  "photoUrls": [
                    "Ve3Eyk4pvHUm-fHV"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 99,
                      "name": "VC2KwW5ZelKSw4qZ OR SLEEP(4)=0 LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 803,
                    "name": "r0FD0OjK9Sao_gen"
                  },
                  "id": 875,
                  "name": "doggie",
                  "photoUrls": [
                    "u7sckbZHfrFSYexu"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 450,
                      "name": "aXWI-ARDf8UKmI3C') OR SLEEP(4)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 297,
                    "name": "9YIeM8DcqbX1tNX3"
                  },
                  "id": 125,
                  "name": "doggie",
                  "photoUrls": [
                    "0GHPJbv71l1KwcP4"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 537,
                      "name": "8YuHu9dV5Fxd9yEB';select pg_sleep(2); -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 693,
                    "name": "grXSKpiJClZFz-Pp"
                  },
                  "id": 638,
                  "name": "doggie",
                  "photoUrls": [
                    "hrj3bnTk5UxdxCnE"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 832,
                      "name": "au9mNsWI6zzjgLmf'));select pg_sleep(8); -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 875,
                    "name": "lWsZ-bhXs4PjBp4L"
                  },
                  "id": 97,
                  "name": "doggie",
                  "photoUrls": [
                    "1_EGGiCBqibs-WOd"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 801,
                      "name": "jGDGDKBx9TFO-2WW') AND SLEEP(2)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 938,
                    "name": "UMmw0YG6mZmPr7AX"
                  },
                  "id": 95,
                  "name": "doggie",
                  "photoUrls": [
                    "kU6tkciKW5DRqApV"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 932,
                      "name": "b-nz6TZQX0F-t3Or'));WAITFOR DELAY '00:00:2'-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 439,
                    "name": "2ZMo93NoMSJqMPvH"
                  },
                  "id": 592,
                  "name": "doggie",
                  "photoUrls": [
                    "uJGQvDrZTRPSwv1f"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 671,
                      "name": "_Yiyo9oQt6E6DWRL);WAITFOR DELAY '00:00:8'-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 369,
                    "name": "zceQvuKSxetHO50Y"
                  },
                  "id": 1001,
                  "name": "doggie",
                  "photoUrls": [
                    "gRs5wB99ZFy_pg9F"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 798,
                      "name": "zQLimzJbR_zZQsEo') AND SLEEP(8)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 464,
                    "name": "6KIaSMkZvARbBI7T"
                  },
                  "id": 935,
                  "name": "doggie",
                  "photoUrls": [
                    "5s0nFNVNW2QEHj2n"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 915,
                      "name": "Qwn97DzLWTrJcah6& sleep 8"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 487,
                    "name": "gL-QewFEKBelWxe-"
                  },
                  "id": 463,
                  "name": "doggie",
                  "photoUrls": [
                    "h8kf2I6MTEaVB0mx"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 154,
                      "name": "rjUwczwDi8Sh1JR6');WAITFOR DELAY '00:00:8'-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 524,
                    "name": "twnKqD8wE5tQZITM"
                  },
                  "id": 677,
                  "name": "doggie",
                  "photoUrls": [
                    "flfHZjgkKdRwoeMh"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 382,
                      "name": "pJRUaJV2Vtc4nNEa'));WAITFOR DELAY '00:00:2'-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 67,
                    "name": "prXArLTTD9HqFRzG"
                  },
                  "id": 813,
                  "name": "doggie",
                  "photoUrls": [
                    "99CU2ackaTCiFKdN"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 187,
                      "name": "6O256gyH_0oUMWGD) AND SLEEP(4)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 667,
                    "name": "wK-gZ07-4Y3LlS2s"
                  },
                  "id": 699,
                  "name": "doggie",
                  "photoUrls": [
                    "dQmYCkqS9RkW-m9G"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 1024,
                      "name": "zQUA2VNoLnHIgAKy') AND SLEEP(4)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 21,
                    "name": "bEikrLlRFKCJIs4k"
                  },
                  "id": 536,
                  "name": "doggie",
                  "photoUrls": [
                    "UNEXBG_cPJ1FMVpV"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 885,
                      "name": "VGKZV_9GhgiQdiS8' AND 1 IN (SELECT BENCHMARK(2*15000000,MD5(CHAR(97))))-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 984,
                    "name": "RE3sFxwU03lQUeY2"
                  },
                  "id": 574,
                  "name": "doggie",
                  "photoUrls": [
                    "Gg_q59RjsS3Eee3V"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 128,
                      "name": "tPYt6fGLwCb5KadA' AND 1 IN (SELECT BENCHMARK(2*15000000,MD5(CHAR(97))))-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 281,
                    "name": "7c1H6y3SxzbQ8kd0"
                  },
                  "id": 583,
                  "name": "doggie",
                  "photoUrls": [
                    "fZP9br8-FlzMEtKi"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 103,
                      "name": "_jSU5HCbb9dWH2uv';WAITFOR DELAY '00:00:8'-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 325,
                    "name": "_v1e5vJcGtThE8vv"
                  },
                  "id": 778,
                  "name": "doggie",
                  "photoUrls": [
                    "aSYTNTi4_9xqyn2s"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 106,
                      "name": "gC7MU9A5uB7NYXft=SLEEP(2)="
                    }
                  ]
                },
                {
                  "category": {
                    "id": 411,
                    "name": "WpGvNeEMM5TrQjS9"
                  },
                  "id": 820,
                  "name": "doggie",
                  "photoUrls": [
                    "PUjNahrl07t2Le-W"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 277,
                      "name": "SNseIGShZib7HxnO'));WAITFOR DELAY '00:00:8'-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 312,
                    "name": "PH9NisHcDbRzloMu"
                  },
                  "id": 332,
                  "name": "doggie",
                  "photoUrls": [
                    "byWFug6UgzW9QgDX"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 912,
                      "name": "sozcZ6a4lwTwtyA5') AND SLEEP(8)=0 LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 155,
                    "name": "Gin7DkMhbLjgNzoh"
                  },
                  "id": 605,
                  "name": "doggie",
                  "photoUrls": [
                    "GWFiO24kasAC7LgJ"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 297,
                      "name": "XMc_gT_x-CSCyeKt"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 995,
                    "name": "YVLaZ5YgIVz2mzem"
                  },
                  "id": 824,
                  "name": "doggie",
                  "photoUrls": [
                    "cw8RIR9_bsh_K5_X"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 888,
                      "name": "XnDukC32EJZgtWDH') AND SLEEP(2)=0 LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 942,
                    "name": "_3xH4ExwxwdcDWwA"
                  },
                  "id": 929,
                  "name": "doggie",
                  "photoUrls": [
                    "wfZW9NQ2TOuABiVC"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 991,
                      "name": "-z3ZMI2HkmLElVyp' WHERE SLEEP(8) LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 401,
                    "name": "KNkEmyt9KHjq_AY9"
                  },
                  "id": 247,
                  "name": "doggie",
                  "photoUrls": [
                    "MsvrY536obSDniji"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 508,
                      "name": "-WejlqZV1u6IH5mx' WHERE SLEEP(4) LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 935,
                    "name": "Yzsgzo5z8jej-3nr"
                  },
                  "id": 581,
                  "name": "doggie",
                  "photoUrls": [
                    "FT-MEVb-Q3LOU2Ct"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 14,
                      "name": "lHriamoz1Tkb4mx1 WHERE SLEEP(4) LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 451,
                    "name": "gxS3BbhoaYnj4JvI"
                  },
                  "id": 611,
                  "name": "doggie",
                  "photoUrls": [
                    "7-VZJ2185DG0g8oX"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 841,
                      "name": "IHYdE8dcp90Rfk7Z' OR SLEEP(8)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 843,
                    "name": "lGiEi-Px4yZ_s7DJ"
                  },
                  "id": 302,
                  "name": "doggie",
                  "photoUrls": [
                    "wUWB7ac23Ose29AJ"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 515,
                      "name": "PX6Put5L5UdSo6Mh') AND SLEEP(8)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 976,
                    "name": "YoY5yz5UNWUyPbU5"
                  },
                  "id": 967,
                  "name": "doggie",
                  "photoUrls": [
                    "9vr1QsY-gVAVrFvQ"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 201,
                      "name": "OZxec92og3J9UIku)) OR SLEEP(4)=0 LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 441,
                    "name": "9gMxx3fSwGsPHYO3"
                  },
                  "id": 471,
                  "name": "doggie",
                  "photoUrls": [
                    "amnjgHHmsr0iO_li"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 134,
                      "name": "g5ckNfUIvL806F6F' WHERE SLEEP(4) LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 738,
                    "name": "1VCEMiErDUuwuwBE"
                  },
                  "id": 649,
                  "name": "doggie",
                  "photoUrls": [
                    "Conq-5vfGG3rTgsu"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 922,
                      "name": "kv0qTb0Zefo1LbKM));select pg_sleep(8); -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 443,
                    "name": "ytY_lCozktMzkGrP"
                  },
                  "id": 711,
                  "name": "doggie",
                  "photoUrls": [
                    "YmtdIcFRLf2Ci5Ug"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 947,
                      "name": "hZVHdUzrm9K0eJ_5' OR SLEEP(4)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 235,
                    "name": "wxVz1C-QYLxNWK7H"
                  },
                  "id": 953,
                  "name": "doggie",
                  "photoUrls": [
                    "j_BIEKy9bAeCTFwE"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 767,
                      "name": "ZOb5iRFCOpLMMl8z') AND SLEEP(2)=0 LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 516,
                    "name": "O0hWcPHs0jexT4Zm"
                  },
                  "id": 568,
                  "name": "doggie",
                  "photoUrls": [
                    "tXzLp1jZES9pdKyl"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 341,
                      "name": "DoCSPKZ-DXes2YU8') OR SLEEP(4)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 852,
                    "name": "_04hMQj9o2mZT8C7"
                  },
                  "id": 979,
                  "name": "doggie",
                  "photoUrls": [
                    "o2_k8cUvb9nQxhJT"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 458,
                      "name": "4AAAf3qfew_C4puc"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 984,
                    "name": "wS4nAODvzwApixca"
                  },
                  "id": 376,
                  "name": "doggie",
                  "photoUrls": [
                    "Qq212YG5f_aCSnd2"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 248,
                      "name": "3JrUBNhAXCYUNgiI OR SLEEP(4)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 423,
                    "name": "5sHEerD8Td75spq0"
                  },
                  "id": 195,
                  "name": "doggie",
                  "photoUrls": [
                    "7VCcG6A5uaZ7-Wnq"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 772,
                      "name": "bhJLBPc-HE2sS33T') OR SLEEP(4)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 975,
                    "name": "m_RSCD6tEBQn0Xzn"
                  },
                  "id": 532,
                  "name": "doggie",
                  "photoUrls": [
                    "b02t0Vzxni1KKnGI"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 658,
                      "name": "cCvg3q_azccIEwD- WHERE SLEEP(8) LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 497,
                    "name": "5wWgyiDPpGy-a0iy"
                  },
                  "id": 1015,
                  "name": "doggie",
                  "photoUrls": [
                    "v-PKgH8Ucq2FKT92"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 88,
                      "name": "fAQVpVVnoD6Mk3Ck));WAITFOR DELAY '00:00:2'-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 11,
                    "name": "94vXZv5EM14Mxxwg"
                  },
                  "id": 1017,
                  "name": "doggie",
                  "photoUrls": [
                    "fM0mhnD9zVihWnYM"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 333,
                      "name": "Slgp5jlpEBUYZN1O)) OR SLEEP(4)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 301,
                    "name": "lbKzNHGu4pqPkApx"
                  },
                  "id": 432,
                  "name": "doggie",
                  "photoUrls": [
                    "I1-gjwO6orXM0rrB"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 178,
                      "name": "tnOdNoLBuVdK51ZM) AND SLEEP(2)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 775,
                    "name": "1gVpNO1CMEJbHIgI"
                  },
                  "id": 783,
                  "name": "doggie",
                  "photoUrls": [
                    "jGsBppkvi6zJa9xg"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 964,
                      "name": "DSwWSkA31_pMEBh_)) AND SLEEP(4)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 635,
                    "name": "vM4f39HLuGzf3t5U"
                  },
                  "id": 712,
                  "name": "doggie",
                  "photoUrls": [
                    "c73E4TJMfbL2q-Tl"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 221,
                      "name": "d-wPSIe-7HlBXizQ') AND SLEEP(2)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 424,
                    "name": "tf3qiiIHvBg4lsbY"
                  },
                  "id": 595,
                  "name": "doggie",
                  "photoUrls": [
                    "AIVmIgnb-A8aJKdC"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 912,
                      "name": "4lO1CtHI9ZKqA1Pd' AND SLEEP(8)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 231,
                    "name": "XLtIM0Q31322IJPO"
                  },
                  "id": 431,
                  "name": "doggie",
                  "photoUrls": [
                    "HnJ6GzyQ_56Kr5LW"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 481,
                      "name": "PrFq78cT6P5NnZbz;select pg_sleep(8); -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 42,
                    "name": "dUwAfT5JzkbSIc0N"
                  },
                  "id": 606,
                  "name": "doggie",
                  "photoUrls": [
                    "oCxJ1YLEHyaxj-tf"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 209,
                      "name": "_H9jeF3vwJ5lhfCB);select pg_sleep(4); -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 827,
                    "name": "y95gSVXt__pWpp0s"
                  },
                  "id": 132,
                  "name": "doggie",
                  "photoUrls": [
                    "CSKZrw7JSygTIUSK"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 309,
                      "name": "fQUa4R9clW9X3m3u';WAITFOR DELAY '00:00:4'-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 0,
                    "name": "string"
                  },
                  "id": 23,
                  "name": "Bobby",
                  "photoUrls": [
                    "string"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 0,
                      "name": "string"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 941,
                    "name": "efHDJgfvcZP0dey1"
                  },
                  "id": 117,
                  "name": "doggie",
                  "photoUrls": [
                    "7Ite9w8rXAIwJMsg"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 170,
                      "name": "JJ8FfSG2WnKx3pZx';select pg_sleep(8); -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 100,
                    "name": "w391IXpoH3KKz5lZ"
                  },
                  "id": 912,
                  "name": "doggie",
                  "photoUrls": [
                    "wemNKCVDqboZebYb"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 675,
                      "name": "file://WEB-INF/web.xml"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 171,
                    "name": "4RWW6R8uoJwe-Msn"
                  },
                  "id": 14,
                  "name": "doggie",
                  "photoUrls": [
                    "XoNR_WunQh-dO08T"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 40,
                      "name": "PY1bWuSXBntDmGXh;WAITFOR DELAY '00:00:8'-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 959,
                    "name": "-HOOJ3vu1wLf0krY"
                  },
                  "id": 201,
                  "name": "doggie",
                  "photoUrls": [
                    "3hzrUfMeLN8nVnx2"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 258,
                      "name": "L4N22PRhJmEbQBbC;WAITFOR DELAY '00:00:4'-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 40,
                    "name": "RxnaEIj-7Zlil7hA"
                  },
                  "id": 584,
                  "name": "doggie",
                  "photoUrls": [
                    "xyJHpiwUDG90-mKF"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 663,
                      "name": "WJARwvTt8A3xnuz1';WAITFOR DELAY '00:00:2'-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 315,
                    "name": "QkC13ZoEZ7WHeHV2"
                  },
                  "id": 330,
                  "name": "doggie",
                  "photoUrls": [
                    "8CtsCa0rzW6UjLK2"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 919,
                      "name": "_49d-hqw1wLZTVlp') OR SLEEP(8)=0 LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 993,
                    "name": "9-rzefgsP8ogXO4Q"
                  },
                  "id": 1014,
                  "name": "doggie",
                  "photoUrls": [
                    "_7G4jU607LhEU0xb"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 602,
                      "name": "Hw7RbXkvWXm_MdWP');WAITFOR DELAY '00:00:2'-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 398,
                    "name": "DxDilhUX22WTBxxg"
                  },
                  "id": 990,
                  "name": "doggie",
                  "photoUrls": [
                    "ReXlmnlOh6KsLzLx"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 666,
                      "name": "qxf8U-ZEo1HGZAFa'));WAITFOR DELAY '00:00:2'-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 337,
                    "name": "mpjHuO4gNfqyvflj"
                  },
                  "id": 401,
                  "name": "doggie",
                  "photoUrls": [
                    "85Rt8S6-qtilttw8"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 273,
                      "name": "rZ1utzZ_V0nmuC8Q') OR SLEEP(8)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 237,
                    "name": "7sP7b135x5VWBsQZ"
                  },
                  "id": 836,
                  "name": "doggie",
                  "photoUrls": [
                    "zQ7GS-t85nb6kBQx"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 190,
                      "name": "\r\nx-tinfoil-response-splitting: true"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 313,
                    "name": "nb8qAw-lL2OZnQVp"
                  },
                  "id": 123456789987656880,
                  "name": "Otussie",
                  "photoUrls": [
                    "string"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 967,
                      "name": "zAVW-cUUysv3AHkh"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 22,
                    "name": "WurhAutWgbsjRL2E"
                  },
                  "id": 947,
                  "name": "doggie",
                  "photoUrls": [
                    "6jzBOa-Z1iqgUi1B"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 436,
                      "name": "h5HjbZkwyaLqQHGZ) OR SLEEP(4)=0 LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 227,
                    "name": "KoKPQ6EnfME-Hf-4"
                  },
                  "id": 280,
                  "name": "doggie",
                  "photoUrls": [
                    "xE6aSBaqtD35-t3P"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 288,
                      "name": "gkJGv8oxsswHA_Mo sleep 4"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 919,
                    "name": "Om_DMHsmhxf2DPBO"
                  },
                  "id": 344,
                  "name": "doggie",
                  "photoUrls": [
                    "lF_F8k642hkWkbzb"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 827,
                      "name": "RE83HStUbxZlV8z0&& sleep 8"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 564,
                    "name": "Y1UssbNDRkXpW7T4"
                  },
                  "id": 139,
                  "name": "doggie",
                  "photoUrls": [
                    "vKoHAoureCYrRZ2t"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 453,
                      "name": "gnsJhkION6_mMlSm"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 957,
                    "name": "cDbyHLCZX1Uu9_te"
                  },
                  "id": 144,
                  "name": "doggie",
                  "photoUrls": [
                    "5HLbUDyKhYiz34vW"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 726,
                      "name": "WPo8IGuftnQLNfZ4& ping -n 8 localhost"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 527,
                    "name": "2zgEozX17TjG_RIS"
                  },
                  "id": 220,
                  "name": "doggie",
                  "photoUrls": [
                    "WJ-jvnx9k_mwWow5"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 738,
                      "name": "G8Gyj81Iq2WsW3CM&& ping -n 8 localhost"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 276,
                    "name": "zZGW_-oHOtttwSFz"
                  },
                  "id": 394,
                  "name": "doggie",
                  "photoUrls": [
                    "WqPXySa2er3wrUN_"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 731,
                      "name": "I8WjAqI3BmQLe__5| ping -n 8 localhost"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 717,
                    "name": "9pn5eUwFeKUlhmWE"
                  },
                  "id": 182,
                  "name": "doggie",
                  "photoUrls": [
                    "N_xacQffq0vS0tXC"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 524,
                      "name": "file://WEB-INF/web.xml"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 568,
                    "name": "RZP3tmc5a1LMLXYL"
                  },
                  "id": 122,
                  "name": "doggie",
                  "photoUrls": [
                    "TEbQ07J9-TchkHNF"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 42,
                      "name": "BNl96gVC5Vx7vTKX& sleep 2"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 108,
                    "name": "85e6qFsWgvtoMN8n"
                  },
                  "id": 244,
                  "name": "doggie",
                  "photoUrls": [
                    "zRlpnYhAsfff7M0P"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 448,
                      "name": "tnA4YyakLyz0HQdw&& sleep 4"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 840,
                    "name": "OjJ_AKheouaSTenE"
                  },
                  "id": 76,
                  "name": "doggie",
                  "photoUrls": [
                    "J8H5cxs6GvGIbjP6"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 246,
                      "name": "bLePO-NSj6wksncE ping -n 2 localhost"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 424,
                    "name": "IJXA74ygeyZrPw4B"
                  },
                  "id": 312,
                  "name": "doggie",
                  "photoUrls": [
                    "__5VqVtb7Url29rT"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 44,
                      "name": "RtONEw6ftQLZLmOZ;select sleep(2);#"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 948,
                    "name": "hAazC_mZafL48DNw"
                  },
                  "id": 863,
                  "name": "doggie",
                  "photoUrls": [
                    "7RxBh4SqE8uUsk93"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 108,
                      "name": "OPdmKWJPMeKBMHTc;select sleep(4);#"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 737,
                    "name": "XE9PIjPMYt7TEVYt"
                  },
                  "id": 377,
                  "name": "doggie",
                  "photoUrls": [
                    "T0p-hU2Hlzt6EsxH"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 331,
                      "name": "XkjWrJjhLDONKHQ2)) OR SLEEP(4)=0 LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 196,
                    "name": "GNSMOK41w0tI2sTQ"
                  },
                  "id": 620,
                  "name": "doggie",
                  "photoUrls": [
                    "ytb6ZtQTYZgKFims"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 154,
                      "name": "w_bX8Hx9SWPBCmPT sleep 4"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 414,
                    "name": "FUszaf9DUK9pCss4"
                  },
                  "id": 958,
                  "name": "doggie",
                  "photoUrls": [
                    "dVTQiswxnPClGr-y"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 761,
                      "name": "VmdXyweh_jf0Hiyx') AND SLEEP(4)=0 LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 932,
                    "name": "fVtCEINopFFIyXJ8"
                  },
                  "id": 336,
                  "name": "doggie",
                  "photoUrls": [
                    "gDMwWco7wyJOXcJd"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 630,
                      "name": "7GQD2hrtxqPQ9Mpc) OR SLEEP(8)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 736,
                    "name": "x8jhOJgee6P3iIGH"
                  },
                  "id": 633,
                  "name": "doggie",
                  "photoUrls": [
                    "5EnU8TiZl_eARsTm"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 991,
                      "name": "VzL4XIR1pfTfFDjR OR SLEEP(4)=0 LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 249,
                    "name": "5I3hz2eFOR2brfR-"
                  },
                  "id": 355,
                  "name": "doggie",
                  "photoUrls": [
                    "TlzFxHrwupjfd2rI"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 939,
                      "name": "qbehHhXdg89eFpdl)) OR SLEEP(4)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 464,
                    "name": "jNw6MjmJJnqmU1Lt"
                  },
                  "id": 433,
                  "name": "doggie",
                  "photoUrls": [
                    "66ncY63It3E5YacU"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 321,
                      "name": "AX-qb-odyTcsOCxO));WAITFOR DELAY '00:00:8'-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 800,
                    "name": "hPKJ4GIXVQEpfSY6"
                  },
                  "id": 537,
                  "name": "doggie",
                  "photoUrls": [
                    "ORqZbjFjB8u9XoRp"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 892,
                      "name": "mvS_51AvSO8lPmuB)) AND SLEEP(2)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 895,
                    "name": "bVl-g6_DqK1XlVNO"
                  },
                  "id": 58,
                  "name": "doggie",
                  "photoUrls": [
                    "oygSALK6dJp85niX"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 799,
                      "name": "F5ZjWkJByjRkTywM'));WAITFOR DELAY '00:00:8'-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 47,
                    "name": "N7M6bMW7iev4pBkD"
                  },
                  "id": 214,
                  "name": "doggie",
                  "photoUrls": [
                    "mv9zH2i_v7Qdq5qj"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 842,
                      "name": "yBXUyWivi9zWcJqm'));select pg_sleep(4); -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 193,
                    "name": "3AOB_lC5gK76tulE"
                  },
                  "id": 113,
                  "name": "doggie",
                  "photoUrls": [
                    "AkxRmnZtd4DFMIzN"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 368,
                      "name": "wa0qLtH44Ww2km4t;WAITFOR DELAY '00:00:4'-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 87,
                    "name": "1sOh6GgG57oB1Rw6"
                  },
                  "id": 670,
                  "name": "doggie",
                  "photoUrls": [
                    "iWZ7LppLCDFziM8a"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 776,
                      "name": "lyrGmavY2w5Kp-nD';WAITFOR DELAY '00:00:2'-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 416,
                    "name": "HhiN57eCB0S-5icE"
                  },
                  "id": 472,
                  "name": "doggie",
                  "photoUrls": [
                    "smFw7z0780BJFjDb"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 573,
                      "name": "9OaGFAeOhV1T9Xwa'=SLEEP(8)='"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 67,
                    "name": "FvN0_DYBpRE-iNsp"
                  },
                  "id": 817,
                  "name": "doggie",
                  "photoUrls": [
                    "tq08Vd0Pw6uIj_ZS"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 32,
                      "name": "KYHoOtAun-bV3QVs'));WAITFOR DELAY '00:00:4'-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 707,
                    "name": "TlzKRAE6MUUxFQDO"
                  },
                  "id": 978,
                  "name": "doggie",
                  "photoUrls": [
                    "BpXa3P8p1ihGnKPb"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 118,
                      "name": "NcAphmziIaikuyDL'));WAITFOR DELAY '00:00:8'-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 138,
                    "name": "w-XA7liaevCmyhMe"
                  },
                  "id": 940,
                  "name": "doggie",
                  "photoUrls": [
                    "617if20WpUp-x6jb"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 371,
                      "name": "6gVzjmBQn6JKP5ZP);WAITFOR DELAY '00:00:2'-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 722,
                    "name": "oIZiHmVHsvcwzMfS"
                  },
                  "id": 789,
                  "name": "doggie",
                  "photoUrls": [
                    "ffIKtHa03JWbpqlq"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 403,
                      "name": "psKmiE6qj9WO3wZy);WAITFOR DELAY '00:00:4'-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 789,
                    "name": "oVvX57tGxU9X75u2"
                  },
                  "id": 632,
                  "name": "doggie",
                  "photoUrls": [
                    "q2PI5KbaK5Aekw3O"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 245,
                      "name": "meEzhiF6fHGMAqQ2"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 114,
                    "name": "_6YQd8Lt8GYmbWjA"
                  },
                  "id": 752,
                  "name": "doggie",
                  "photoUrls": [
                    "yGaz1EtwPRx8TskU"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 128,
                      "name": "ciR5HYAVbmF97U_h sleep 8"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 148,
                    "name": "NRCQ8YWhduL-fiWs"
                  },
                  "id": 314,
                  "name": "doggie",
                  "photoUrls": [
                    "rLD_aLgBHbXbXqMB"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 588,
                      "name": "Izwu5bWac6LI1UG8| sleep 2"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 642,
                    "name": "OG9HneSWjN92F4wc"
                  },
                  "id": 941,
                  "name": "doggie",
                  "photoUrls": [
                    "V4pL11zHlQ9CQchd"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 302,
                      "name": "5iPLXWoFxC-a4Bcn| sleep 8"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 946,
                    "name": "UGDwVui8PIE6yN1R"
                  },
                  "id": 44,
                  "name": "doggie",
                  "photoUrls": [
                    "JXRZRk9fken8O0v0"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 707,
                      "name": "oBacXx_bKcf42xme ping -n 2 localhost"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 607,
                    "name": "q4NUBjNSRo18Q7I0"
                  },
                  "id": 436,
                  "name": "doggie",
                  "photoUrls": [
                    "Ou461wksdTmMCYPh"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 455,
                      "name": "d0VCQmkCrL-3OJ0x ping -n 8 localhost"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 529,
                    "name": "noh6UYnFIrVEvQ_I"
                  },
                  "id": 245,
                  "name": "doggie",
                  "photoUrls": [
                    "etU9TuX1v86L16Fo"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 430,
                      "name": "dfZAWmS7STqvbmK4 ping -n 4 localhost"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 617,
                    "name": "_ZPt_HdikZgQclFT"
                  },
                  "id": 717,
                  "name": "doggie",
                  "photoUrls": [
                    "-XRF2W-Eet4ay-4r"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 62,
                      "name": "kcRuL7-6mDEacT42&& ping -n 8 localhost"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 189,
                    "name": "pqIj6VlIqQkV0i4L"
                  },
                  "id": 255,
                  "name": "doggie",
                  "photoUrls": [
                    "Hpwn-0xBRN0xTESN"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 237,
                      "name": "2YRLAMXK2XVz_nrn| ping -n 8 localhost"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 688,
                    "name": "YVF1SqFZNggHi_HO"
                  },
                  "id": 140,
                  "name": "doggie",
                  "photoUrls": [
                    "-evGPyIedcjrVpAA"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 142,
                      "name": "yAw-soqYSU2ihs0S)"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 24,
                    "name": "bQ9t6rq3ZhRdpWtT"
                  },
                  "id": 921,
                  "name": "doggie",
                  "photoUrls": [
                    "uIqIfSVwwsODb96t"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 659,
                      "name": "6LcTTU87QU6a5Y5u AND 1 IN (SELECT BENCHMARK(4*15000000,MD5(CHAR(97))))-- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 985,
                    "name": "Ao1m3iEXoPXuVQnk"
                  },
                  "id": 894,
                  "name": "doggie",
                  "photoUrls": [
                    "YR4lXsdC4npFu-EY"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 812,
                      "name": "BpqcMqENYGSV9Tel'=SLEEP(4)='"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 788,
                    "name": "73kn5ZavP7LIFyg3"
                  },
                  "id": 133,
                  "name": "doggie",
                  "photoUrls": [
                    "blzZXiCuM9_7VXW-"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 453,
                      "name": "q7gE6QhmtPTmqKtx'=SLEEP(2)='"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 458,
                    "name": "1naQlRuXV077e1bL"
                  },
                  "id": 769,
                  "name": "doggie",
                  "photoUrls": [
                    "pDmPF5C0GwVmsoqy"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 247,
                      "name": "wIyIgYyoVLvdrF-g'=SLEEP(4)='"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 735,
                    "name": "5xW3u_FSBYPBqX5Z"
                  },
                  "id": 951,
                  "name": "doggie",
                  "photoUrls": [
                    "qChkf9XmYlpg-kKB"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 106,
                      "name": "4s_Q1u0iq3UCQpGn') OR SLEEP(4)=0 LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 252,
                    "name": "uVdM6KkD6GDDiOlF"
                  },
                  "id": 242,
                  "name": "doggie",
                  "photoUrls": [
                    "2WWqd4uFI3r4AiDj"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 911,
                      "name": "kwxUJwjMzl13zS2s') OR SLEEP(8)=0 LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 308,
                    "name": "F6mO61ldm600dGNy"
                  },
                  "id": 341,
                  "name": "doggie",
                  "photoUrls": [
                    "D_LZxSIwNoGhU_t_"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 237,
                      "name": "KGk3HfQ8_XgQa9d6 AND SLEEP(8)=0 LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 406,
                    "name": "7c-eaL2FfIosSwhS"
                  },
                  "id": 507,
                  "name": "doggie",
                  "photoUrls": [
                    "EdainYnetu2qbxEs"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 341,
                      "name": "NLUP3dcsvbadOLyh) AND SLEEP(4)=0 LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 7,
                    "name": "yeg3xFqRxKriwUhr"
                  },
                  "id": 561,
                  "name": "doggie",
                  "photoUrls": [
                    "sisD1LSjLdQEsiZq"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 506,
                      "name": "0BTg-s_-ofpjFQM1' AND SLEEP(8)=0 LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 59,
                    "name": "iMMAIStWxCey_OQh"
                  },
                  "id": 765,
                  "name": "doggie",
                  "photoUrls": [
                    "mvV3r5BcQskBEsJM"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 710,
                      "name": "_tJsXaPJ3xoH3-53)) AND SLEEP(8)=0 LIMIT 1 # "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 380,
                    "name": "eBBnbjmPGEZAxmjD"
                  },
                  "id": 623,
                  "name": "doggie",
                  "photoUrls": [
                    "rgc2i4ULgFuc89Aq"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 815,
                      "name": "8R1KDxDP3iUdQdSD OR SLEEP(4)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 918,
                    "name": "pCgkqAWcFEgT4qtb"
                  },
                  "id": 372,
                  "name": "doggie",
                  "photoUrls": [
                    "1eY-8-vvj48KOU9q"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 138,
                      "name": "kqAMZsWXRkqeYA6o' OR SLEEP(8)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 175,
                    "name": "doJFnqM_uSYoG3c5"
                  },
                  "id": 407,
                  "name": "doggie",
                  "photoUrls": [
                    "0MHmXFCYFFyl8JJx"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 841,
                      "name": "7Z-_v7y_qX6TzhTr') AND SLEEP(8)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 1018,
                    "name": "TbPFcVcl1yXxZ3Hs"
                  },
                  "id": 534,
                  "name": "doggie",
                  "photoUrls": [
                    "OYaJpn-ImskBfA6t"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 990,
                      "name": "AoX0oV6znvr1lA7G' AND SLEEP(8)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 665,
                    "name": "wOZmMiUPpMbRYU46"
                  },
                  "id": 365,
                  "name": "doggie",
                  "photoUrls": [
                    "vfSuMp_6hUO0rwyG"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 875,
                      "name": "b-YgQqpjfDUuvW2O') AND SLEEP(2)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 643,
                    "name": "lPijKplh_ZpZMdkn"
                  },
                  "id": 692,
                  "name": "doggie",
                  "photoUrls": [
                    "DVbb6mZ_JfkNdi55"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 384,
                      "name": "L_UHwur2OFT97d6o)) AND SLEEP(8)=0 LIMIT 1 -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 925,
                    "name": "pkyqiJB_7pgmv16c"
                  },
                  "id": 357,
                  "name": "doggie",
                  "photoUrls": [
                    "N_0orKe-ct92a6Dc"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 407,
                      "name": "NDM78kHkfXPmzbhL;select pg_sleep(2); -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 921,
                    "name": "_ICqPl34DMLVm9af"
                  },
                  "id": 451,
                  "name": "doggie",
                  "photoUrls": [
                    "bMrdcB_lRSQMkQuW"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 859,
                      "name": "-7Im_tt_6ZFvj4JV);select pg_sleep(2); -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 225,
                    "name": "mFFpzPNf7p-Gm4A5"
                  },
                  "id": 183,
                  "name": "doggie",
                  "photoUrls": [
                    "8h9oywdt2GiS5XUB"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 801,
                      "name": "RDcltXJU6yQJsZCS));select pg_sleep(4); -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 71,
                    "name": "2A4ccdTSJKUnmFn1"
                  },
                  "id": 456,
                  "name": "doggie",
                  "photoUrls": [
                    "4MQPI5Y7iPBLRN8V"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 793,
                      "name": "Y1G5m2YPYQ_S9fbo';select pg_sleep(2); -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 467,
                    "name": "iMrUjK94T8WLpAJ9"
                  },
                  "id": 80,
                  "name": "doggie",
                  "photoUrls": [
                    "w14Nk5q6LyzlvmXE"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 662,
                      "name": "e-1nKAtIbRThObth'));select pg_sleep(8); -- "
                    }
                  ]
                },
                {
                  "category": {
                    "id": 313,
                    "name": "nb8qAw-lL2OZnQVp"
                  },
                  "id": 123456789987657060,
                  "name": "Otussie",
                  "photoUrls": [
                    "string"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 967,
                      "name": "zAVW-cUUysv3AHkh"
                    }
                  ]
                },
                {
                  "category": {
                    "id": 313,
                    "name": "nb8qAw-lL2OZnQVp"
                  },
                  "id": 123456789987657180,
                  "name": "Otussie",
                  "photoUrls": [
                    "string"
                  ],
                  "status": "sold",
                  "tags": [
                    {
                      "id": 967,
                      "name": "zAVW-cUUysv3AHkh"
                    }
                  ]
                }
              ],
              "_A.InvokeRESTService.responseCodes": null,
              "_A.InvokeRESTService.responseCodesSchema": null,
              "_A.InvokeRESTService.responseType": "",
              "_A.InvokeRESTService.statusCode": 200
            },
            "tasks": {
              "InvokeRESTService": {
                "change": 2,
                "status": 40,
                "input": {
                  "Method": "GET",
                  "Server Certificate": "",
                  "Timeout": 0,
                  "Uri": "https://petstore.swagger.io/v2/pet/findByStatus",
                  "Use certificate for verification": false,
                  "enableASR": false,
                  "headers": {
                    "Connection": "close"
                  },
                  "host": "",
                  "proxy": "",
                  "queryParams": {
                    "status": "sold"
                  },
                  "requestType": "application/json",
                  "resourcePath": "",
                  "serviceName": ""
                }
              },
              "InvokeRESTService1": {
                "change": 1,
                "status": 20,
                "input": null
              }
            },
            "links": {
              "0": {
                "change": 2,
                "status": 0,
                "from": "",
                "to": ""
              },
              "1": {
                "change": 1,
                "status": 2,
                "from": "InvokeRESTService",
                "to": "InvokeRESTService1"
              }
            },
            "returnData": null
          }
        },
        "queueChanges": {
          "2": {
            "change": 2,
            "subflowId": 0,
            "taskId": "InvokeRESTService"
          },
          "3": {
            "change": 0,
            "subflowId": 0,
            "taskId": "InvokeRESTService1"
          }
        }
      },
      {
        "id": 3,
        "flowId": "2e2842d9a30937083491f055d89fa1ec",
        "flowChanges": {
          "0": {
            "newFlow": false,
            "flowURI": "",
            "subflowId": 0,
            "taskId": "InvokeRESTService1",
            "status": 0,
            "attrs": {
              "_A.InvokeRESTService1.configureResponseCodes": false,
              "_A.InvokeRESTService1.headers": {
                "Connection": "keep-alive",
                "Content-Type": "application/json"
              },
              "_A.InvokeRESTService1.responseBody": {
                "category": {
                  "id": 4,
                  "name": "1HU9zwcdiX-17u7f"
                },
                "id": 904,
                "name": "doggie",
                "photoUrls": [
                  "w6PYmahm0YX3NRBu"
                ],
                "status": "sold",
                "tags": [
                  {
                    "id": 449,
                    "name": "xSSey2yoF7az3M04 AND 1 IN (SELECT BENCHMARK(4*15000000,MD5(CHAR(97))))-- "
                  }
                ]
              },
              "_A.InvokeRESTService1.responseCodes": null,
              "_A.InvokeRESTService1.responseCodesSchema": null,
              "_A.InvokeRESTService1.responseType": "",
              "_A.InvokeRESTService1.statusCode": 200
            },
            "tasks": {
              "InvokeRESTService1": {
                "change": 2,
                "status": 40,
                "input": {
                  "Method": "GET",
                  "Server Certificate": "",
                  "Timeout": 0,
                  "Uri": "https://petstore.swagger.io/v2/pet/{id}",
                  "Use certificate for verification": false,
                  "enableASR": false,
                  "host": "",
                  "pathParams": {
                    "id": "904"
                  },
                  "proxy": "",
                  "requestType": "application/json",
                  "resourcePath": "",
                  "serviceName": ""
                }
              },
              "Log200": {
                "change": 1,
                "status": 20,
                "input": null
              },
              "Log405": {
                "change": 2,
                "status": 50,
                "input": null
              }
            },
            "links": {
              "1": {
                "change": 2,
                "status": 0,
                "from": "",
                "to": ""
              },
              "2": {
                "change": 1,
                "status": 2,
                "from": "InvokeRESTService1",
                "to": "Log200"
              },
              "6": {
                "change": 2,
                "status": 1,
                "from": "InvokeRESTService1",
                "to": "Log405"
              }
            },
            "returnData": null
          }
        },
        "queueChanges": {
          "3": {
            "change": 2,
            "subflowId": 0,
            "taskId": "InvokeRESTService1"
          },
          "4": {
            "change": 0,
            "subflowId": 0,
            "taskId": "Log200"
          }
        }
      },
      {
        "id": 4,
        "flowId": "2e2842d9a30937083491f055d89fa1ec",
        "flowChanges": {
          "0": {
            "newFlow": false,
            "flowURI": "",
            "subflowId": 0,
            "taskId": "Log200",
            "status": 0,
            "attrs": null,
            "tasks": {
              "Log200": {
                "change": 2,
                "status": 40,
                "input": {
                  "Log Level": "INFO",
                  "flowInfo": false,
                  "message": "In 200 status"
                }
              },
              "Sleep": {
                "change": 1,
                "status": 20,
                "input": null
              }
            },
            "links": {
              "2": {
                "change": 2,
                "status": 0,
                "from": "",
                "to": ""
              },
              "3": {
                "change": 1,
                "status": 2,
                "from": "Log200",
                "to": "Sleep"
              }
            },
            "returnData": null
          }
        },
        "queueChanges": {
          "4": {
            "change": 2,
            "subflowId": 0,
            "taskId": "Log200"
          },
          "5": {
            "change": 0,
            "subflowId": 0,
            "taskId": "Sleep"
          }
        }
      },
      {
        "id": 5,
        "flowId": "2e2842d9a30937083491f055d89fa1ec",
        "flowChanges": {
          "0": {
            "newFlow": false,
            "flowURI": "",
            "subflowId": 0,
            "taskId": "Sleep",
            "status": 0,
            "attrs": null,
            "tasks": {
              "Sleep": {
                "change": 2,
                "status": 40,
                "input": {
                  "Interval": 0,
                  "IntervalSetting": 3,
                  "IntervalTypeSetting": "Second"
                }
              },
              "Sleep1": {
                "change": 1,
                "status": 20,
                "input": null
              }
            },
            "links": {
              "3": {
                "change": 2,
                "status": 0,
                "from": "",
                "to": ""
              },
              "4": {
                "change": 1,
                "status": 2,
                "from": "Sleep",
                "to": "Sleep1"
              }
            },
            "returnData": null
          }
        },
        "queueChanges": {
          "5": {
            "change": 2,
            "subflowId": 0,
            "taskId": "Sleep"
          },
          "6": {
            "change": 0,
            "subflowId": 0,
            "taskId": "Sleep1"
          }
        }
      },
      {
        "id": 6,
        "flowId": "2e2842d9a30937083491f055d89fa1ec",
        "flowChanges": {
          "0": {
            "newFlow": false,
            "flowURI": "",
            "subflowId": 0,
            "taskId": "Sleep1",
            "status": 0,
            "attrs": null,
            "tasks": {
              "Return": {
                "change": 1,
                "status": 20,
                "input": null
              },
              "Sleep1": {
                "change": 2,
                "status": 40,
                "input": {
                  "Interval": 0,
                  "IntervalSetting": 2,
                  "IntervalTypeSetting": "Second"
                }
              }
            },
            "links": {
              "4": {
                "change": 2,
                "status": 0,
                "from": "",
                "to": ""
              },
              "5": {
                "change": 1,
                "status": 2,
                "from": "Sleep1",
                "to": "Return"
              }
            },
            "returnData": null
          }
        },
        "queueChanges": {
          "6": {
            "change": 2,
            "subflowId": 0,
            "taskId": "Sleep1"
          },
          "7": {
            "change": 0,
            "subflowId": 0,
            "taskId": "Return"
          }
        }
      },
      {
        "id": 7,
        "flowId": "2e2842d9a30937083491f055d89fa1ec",
        "flowChanges": {
          "0": {
            "newFlow": false,
            "flowURI": "res://flow:Test2",
            "subflowId": 0,
            "taskId": "Return",
            "status": 500,
            "attrs": null,
            "tasks": {
              "Return": {
                "change": 2,
                "status": 40,
                "input": null
              }
            },
            "links": {
              "5": {
                "change": 2,
                "status": 0,
                "from": "",
                "to": ""
              }
            },
            "returnData": {
              "output1": "",
              "output2": "asdasdasd"
            }
          }
        },
        "queueChanges": {
          "7": {
            "change": 2,
            "subflowId": 0,
            "taskId": "Return"
          }
        }
      }
    ]`

func TestTaskSubflowErro22r(t *testing.T) {
	var steps []*state.Step
	err := json.Unmarshal([]byte(testData), &steps)
	if err != nil {
		t.Fatal(err)
	}

	var tasks []*Task
	for _, s := range steps {
		tt, err := StepToTask(s)
		if err != nil {
			t.Fatal(err)
		}
		tasks = append(tasks, tt...)

	}
	v, _ := json.Marshal(tasks)
	fmt.Println(string(v))
}

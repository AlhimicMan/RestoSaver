var FrontRest = function(id, name, imgUrl, description, tags) {
    this.ID = id
    this.ImageURL = imgUrl
    this.Name = name
    this.DescrFull = description
    if ((description != undefined) &&(description.length > 160)) {
        let removeIndex = 155
        for (let i = removeIndex; i > removeIndex-20; i--) {
            if (description[i] == " ") {
                removeIndex = i
                break
            }
        }
        this.Descr = description.substr(0, removeIndex) + " ..."
    } else {
        this.Descr = description
    }
    this.Tags = []
    if (Array.isArray(tags)) {
        this.Tags = tags
    }
    
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

var app = new Vue({
    el: '#app',
    data: {
    message: 'RestoSaver',
    MaxPlacesPerPage: 9,
    ShowOnPageNow: 9,
    NewPlace: {
        Name: "",
        Descr: "",
        Web: "",
        Address: "",
        Email: "",
        Phone: "",
    },
    NewPlaceValidator: {
        NameInvalid: false,
        WebInvalid: false,
        AddressInvalid: false,
        EmailInvalid: false,
        PhoneInvalid: false,
    },
    NewGCard: {
        Name: "",
        Email: "",
        Phone: "",
        Summ: "",
    },
    NewGCardValidator: {
        NameInvalid: false,
        EmailInvalid: false,
        PhoneInvalid: false,
        SummInvalid: false,
    },
    placesLoaded: false,
    curPlace : new FrontRest,
    placeSent: false,
    reqSent: false,
    restsListAll: [], //To store all data from backend
    restsListView: [],
    SearchQuery: "",
    },
    methods: {
        SendNewPlace() {
            const keys = Object.keys(this.NewPlaceValidator)
            for (const key of keys) {
                this.NewPlaceValidator[key] = false
            }
            let haveErors = false
            console.log("Saving place")
            if (this.NewPlace.Name === "") {
                this.NewPlaceValidator.NameInvalid = true
                haveErors = true
            }
            if (this.NewPlace.Web === "") {
                this.NewPlaceValidator.WebInvalid = true
                haveErors = true
            }
            if (this.NewPlace.Address === "") {
                this.NewPlaceValidator.AddressInvalid = true
                haveErors = true
            }
            if (this.NewPlace.Email === "") {
                this.NewPlaceValidator.EmailInvalid = true
                haveErors = true
            } else {
                let reEmail = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
                if (!reEmail.test(this.NewPlace.Email)) {
                    this.NewPlaceValidator.EmailInvalid = true
                    haveErors = true
                }
            }
            if (this.NewPlace.Phone === "") {
                this.NewPlaceValidator.PhoneInvalid = true
                haveErors = true
            }
            if (!haveErors) {
                let vApp = this
                axios.post('/api/rests/add', this.NewPlace)
                  .then(async function (response) {
                      console.log(response)
                    if (response.statusText == "OK") {
                        vApp.placeSent = true
                        await sleep (4000);
                        $("#addNewPlace").modal('hide');
                        vApp.placeSent = false
                    }
                  })
                  .catch(function (error) {
                    console.log(error);
                  });
            }
        },
        ShowPlaceCard(place) {
            this.curPlace = place
        },
        SendGCardReq() {
            let buyReqCard = {
                Name: "",
                Email: "",
                Phone: "",
                Summ: 0,
            }
            const keys = Object.keys(this.NewGCardValidator)
            for (const key of keys) {
                this.NewGCardValidator[key] = false
            }
            let haveErors = false
            if (this.NewGCard.Name == "") {
                this.NewGCardValidator.NameInvalid = true
                haveErors = true
            }
            if (this.NewGCard.Email == "") {
                this.NewGCardValidator.EmailInvalid = true
                haveErors = true
            } else {
                let reEmail = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
                if (!reEmail.test(this.NewGCard.Email)) {
                    this.NewGCardValidator.EmailInvalid = true
                    haveErors = true
                }
            }
            if (this.NewGCard.Phone == "") {
                this.NewGCardValidator.PhoneInvalid = true
                haveErors = true
            } else {
                let phoneRe = /^[0-9+-]{9,17}$/;
                if (!phoneRe.test(this.NewGCard.Phone)) {
                    this.NewGCardValidator.PhoneInvalid = true
                    haveErors = true
                }
            }
            if (this.NewGCard.Summ == 0) {
                this.NewGCardValidator.SummInvalid = true
                haveErors = true
            } else {
                let sum = parseInt(this.NewGCard.Summ)
                if ((sum == NaN) || (sum < 500)) {
                    this.NewGCardValidator.SummInvalid = true
                    haveErors = true
                } else {
                    buyReqCard.Summ = sum
                }
            }
            
            if (!haveErors) {
                buyReqCard.Name = this.NewGCard.Name
                buyReqCard.Email = this.NewGCard.Email
                buyReqCard.Phone = this.NewGCard.Phone
                let vApp = this
                axios.post('/api/card/buy', buyReqCard)
                  .then(async function (response) {
                      console.log(response)
                    if (response.statusText == "OK") {
                        vApp.reqSent = true
                        await sleep (2000);
                        $("#placeCard").modal('hide');
                        vApp.reqSent = false
                    }
                  })
                  .catch(function (error) {
                    console.log(error);
                  });
            }
        },
        ShowMorePlaces() {
            console.log("sh more")
            this.ShowOnPageNow += this.MaxPlacesPerPage
            if (this.ShowOnPageNow > this.restsListView.length) {
                this.ShowOnPageNow = this.restsListView.length
            }
        }
    },
    computed: {
        restsList: function() {
            let tempStore = []
            let curMaxVal = this.ShowOnPageNow
            if (this.restsListView.length < curMaxVal) {
                curMaxVal = this.restsListView.length
            }
            for (let i = 0; i < curMaxVal; i++) {
                tempStore.push(this.restsListView[i])
            }
            return tempStore
        }
    },
    watch: {
        SearchQuery(newQ) {
            this.ShowOnPageNow = this.MaxPlacesPerPage
            let nRestList = []
            this.restsListAll.forEach(rest => {
                let nameLow = rest.Name.toLowerCase()
                let searchQLow = newQ.toLowerCase()
                if (nameLow.includes(searchQLow)) {
                    nRestList.push(rest)
                }
            });
            this.restsListView = nRestList
        }
    },
    mounted() {
    //Get list of restaurants
    var isshow = localStorage.getItem('isshow');
    if (isshow== null) {
        localStorage.setItem('isshow', 1);

        // Show popup here
        $("#aboutUs").modal('show');
    }
    let vApp = this
    this.ShowOnPageNow = this.MaxPlacesPerPage
    axios.get('/api/rests/all')
    .then(async function (response) {
        // handle success
        vApp.placesLoaded = true
        if (!Array.isArray(response.data)) {
            console.log("Error getting data")
        }
        response.data.forEach(rest => {
            let rRest = new FrontRest(rest.ID, rest.Name, rest.ImageURL, rest.Descr, rest.Tags)
            vApp.restsListAll.push(rRest)
        });
        console.log("Loaded items: ", vApp.restsListAll.length)
        vApp.restsListView = vApp.restsListAll
    })
    .catch(function (error) {
        // handle error
        console.log(error);
    })

    $(window).scroll(function() {
        // console.log("----------")
        // console.log($(window).scrollTop())
        // console.log($(window).height())
        // console.log($(document).height())
        // console.log("----------")
        if($(window).scrollTop() + $(window).height() > $(document).height() - 100) {
            vApp.ShowMorePlaces()
        }
     });
    }
  })
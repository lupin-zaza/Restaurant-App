﻿using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Identity.EntityFrameworkCore;

namespace Restaurant.Server.Models
{
    public class User : IdentityUser
    {
        public virtual  ICollection<Order> Orders { get; set; }

        public UserProfile UserProfile { get; set; }
    }
}
